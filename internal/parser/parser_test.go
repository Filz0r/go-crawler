package parser

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := GetHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := GetFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLComplex(t *testing.T) {
	tests := []struct {
		name      string
		inputHTML string
		expected  string
	}{
		{
			name:      "h1 preferred over h2",
			inputHTML: "<html><body><h2>Subtitle</h2><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "h2 fallback when no h1",
			inputHTML: "<html><body><h2>Fallback Title</h2></body></html>",
			expected:  "Fallback Title",
		},
		{
			name:      "no heading",
			inputHTML: "<html><body><p>No heading here</p></body></html>",
			expected:  "",
		},
		{
			name:      "empty body",
			inputHTML: "<html><body></body></html>",
			expected:  "",
		},
		{
			name:      "h1 inside nested elements",
			inputHTML: "<html><body><div><section><h1>Nested Title</h1></section></div></body></html>",
			expected:  "Nested Title",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetHeadingFromHTML(tc.inputHTML)
			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: expected %q, got %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLComplex(t *testing.T) {
	tests := []struct {
		name      string
		inputHTML string
		expected  string
	}{
		{
			name:      "main only paragraph",
			inputHTML: `<html><body><main><p>Main paragraph.</p></main></body></html>`,
			expected:  "Main paragraph.",
		},
		{
			name:      "fallback to first paragraph when no main",
			inputHTML: `<html><body><aside><p>Sidebar paragraph.</p></aside></body></html>`,
			expected:  "Sidebar paragraph.",
		},
		{
			name:      "no paragraph",
			inputHTML: "<html><body><h1>Just a heading</h1></body></html>",
			expected:  "",
		},
		{
			name:      "empty body",
			inputHTML: "<html><body></body></html>",
			expected:  "",
		},
		{
			name:      "paragraph inside nested elements",
			inputHTML: "<html><body><div><article><p>Nested paragraph.</p></article></div></body></html>",
			expected:  "Nested paragraph.",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetFirstParagraphFromHTML(tc.inputHTML)
			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: expected %q, got %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := GetURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLComplex(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		inputHTML string
		expected  []string
		wantErr   bool
	}{
		{
			name:      "relative href resolves against base url",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="/blog/post-1">Post</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post-1"},
		},
		{
			name:      "multiple anchor tags found",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="/one">One</a><a href="/two">Two</a><a href="/three">Three</a></body></html>`,
			expected:  []string{"https://crawler-test.com/one", "https://crawler-test.com/two", "https://crawler-test.com/three"},
		},
		{
			name:      "mixed absolute and relative hrefs",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="https://google.com/search">Google</a><a href="/local">Local</a></body></html>`,
			expected:  []string{"https://google.com/search", "https://crawler-test.com/local"},
		},
		{
			name:      "anchor without href is ignored",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a name="not-a-link">No href</a><a href="/real">Real</a></body></html>`,
			expected:  []string{"https://crawler-test.com/real"},
		},
		{
			name:      "no anchors returns empty slice",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><p>No links here</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "anchors nested in other elements are found",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><div><main><a href="/nested">Nested</a></main></div></body></html>`,
			expected:  []string{"https://crawler-test.com/nested"},
		},
		{
			name:      "relative href with query and fragment",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="/blog?page=2#comments">Blog</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog?page=2#comments"},
		},
		{
			name:      "malformed href returns error",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="http://[::1">Bad</a></body></html>`,
			wantErr:   true,
		},
		{
			name:      "relative href with query and fragment",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><a href="https://crawler-test.com/blog?page=2#comments">Blog</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog?page=2#comments"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse base URL: %v", i, tc.name, err)
				return
			}

			actual, err := GetURLsFromHTML(tc.inputHTML, baseURL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Test %v - '%s' FAIL: expected error, got URLs: %v", i, tc.name, actual)
				}
				return
			}
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := GetImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLComplex(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		inputHTML string
		expected  []string
		wantErr   bool
	}{
		{
			name:      "multiple img tags found",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><img src="/one.png"><img src="/two.png"><img src="/three.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/one.png", "https://crawler-test.com/two.png", "https://crawler-test.com/three.png"},
		},
		{
			name:      "mixed absolute and relative srcs",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><img src="https://cdn.example.com/external.png"><img src="/local.png"></body></html>`,
			expected:  []string{"https://cdn.example.com/external.png", "https://crawler-test.com/local.png"},
		},
		{
			name:      "img without src is ignored",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><img alt="no src"><img src="/real.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/real.png"},
		},
		{
			name:      "no images returns empty slice",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><p>No images here</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "images nested in other elements are found",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><div><main><img src="/nested.png"></main></div></body></html>`,
			expected:  []string{"https://crawler-test.com/nested.png"},
		},
		{
			name:      "relative src with query string",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><img src="/thumb.png?w=200&h=200"></body></html>`,
			expected:  []string{"https://crawler-test.com/thumb.png?w=200&h=200"},
		},
		{
			name:      "url with fragment and no path",
			baseURL:   "https://crawler-test.com/blog",
			inputHTML: `<html><body><img src="#local-anchor"></body></html>`,
			expected:  []string{"https://crawler-test.com/blog#local-anchor"},
		},
		{
			name:      "malformed src returns error",
			baseURL:   "https://crawler-test.com",
			inputHTML: `<html><body><img src="http://[::1"></body></html>`,
			wantErr:   true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse base URL: %v", i, tc.name, err)
				return
			}

			actual, err := GetImagesFromHTML(tc.inputHTML, baseURL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Test %v - '%s' FAIL: expected error, got images: %v", i, tc.name, actual)
				}
				return
			}
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataFullPage(t *testing.T) {
	inputBody := `<html><body>
		<h1>Go Crawler</h1>
		<main>
			<p>A crawler that scrapes the web.</p>
		</main>
		<a href="/docs">Docs</a>
		<a href="https://google.com">Google</a>
		<img src="/logo.png">
		<img src="https://cdn.example.com/banner.png">
	</body></html>`

	actual := extractPageData(inputBody, "https://crawler-test.com")

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Go Crawler",
		FirstParagraph: "A crawler that scrapes the web.",
		OutgoingLinks:  []string{"https://crawler-test.com/docs", "https://google.com"},
		ImageURLs:      []string{"https://crawler-test.com/logo.png", "https://cdn.example.com/banner.png"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataH2Fallback(t *testing.T) {
	inputBody := `<html><body>
		<h2>Fallback Heading</h2>
		<p>Some paragraph.</p>
	</body></html>`

	actual := extractPageData(inputBody, "https://crawler-test.com/fallback")

	expected := PageData{
		URL:            "https://crawler-test.com/fallback",
		Heading:        "Fallback Heading",
		FirstParagraph: "Some paragraph.",
		OutgoingLinks:  []string{},
		ImageURLs:      []string{},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataMainParagraphPriority(t *testing.T) {
	inputBody := `<html><body>
		<h1>Main Priority Page</h1>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`

	actual := extractPageData(inputBody, "https://crawler-test.com/landing")

	expected := PageData{
		URL:            "https://crawler-test.com/landing",
		Heading:        "Main Priority Page",
		FirstParagraph: "Main paragraph.",
		OutgoingLinks:  []string{},
		ImageURLs:      []string{},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataEmptyBody(t *testing.T) {
	actual := extractPageData("<html><body></body></html>", "https://crawler-test.com/empty")

	expected := PageData{
		URL:            "https://crawler-test.com/empty",
		Heading:        "",
		FirstParagraph: "",
		OutgoingLinks:  []string{},
		ImageURLs:      []string{},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataLinksRelativeToPageURL(t *testing.T) {
	inputBody := `<html><body>
		<h1>Docs</h1>
		<a href="next-page">Next</a>
		<a href="/absolute">Absolute</a>
		<img src="images/pic.png">
	</body></html>`

	// trailing slash matters: relative links resolve against the page's directory
	actual := extractPageData(inputBody, "https://crawler-test.com/section/")

	expected := PageData{
		URL:            "https://crawler-test.com/section/",
		Heading:        "Docs",
		FirstParagraph: "",
		OutgoingLinks:  []string{"https://crawler-test.com/section/next-page", "https://crawler-test.com/absolute"},
		ImageURLs:      []string{"https://crawler-test.com/section/images/pic.png"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
