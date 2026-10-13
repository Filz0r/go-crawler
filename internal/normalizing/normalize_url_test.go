package normalizing

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		wantErr  bool
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "test localhost",
			inputURL: "localhost",
			expected: "localhost",
		},
		{
			name:     "remove http scheme",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http and https normalize the same",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "no path",
			inputURL: "https://www.boot.dev",
			expected: "www.boot.dev",
		},
		{
			name:     "root path only",
			inputURL: "https://www.boot.dev/",
			expected: "www.boot.dev/",
		},
		{
			name:     "trailing slash is kept",
			inputURL: "https://www.boot.dev/blog/",
			expected: "www.boot.dev/blog/",
		},
		{
			name:     "query string is dropped",
			inputURL: "https://www.boot.dev/blog/path?sort=asc&page=2",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "fragment is dropped",
			inputURL: "https://www.boot.dev/blog/path#section",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "port is kept",
			inputURL: "http://localhost:8080/health",
			expected: "localhost:8080/health",
		},
		{
			name:     "userinfo is dropped",
			inputURL: "https://user:pass@www.boot.dev/blog",
			expected: "www.boot.dev/blog",
		},
		{
			name:     "ip address host",
			inputURL: "https://127.0.0.1:4000/docs",
			expected: "127.0.0.1:4000/docs",
		},
		{
			name:     "deeply nested path",
			inputURL: "https://www.boot.dev/a/b/c/d/e",
			expected: "www.boot.dev/a/b/c/d/e",
		},
		{
			name:     "invalid URL errors",
			inputURL: "http://[::1",
			wantErr:  true,
		},
		// add more test cases here
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeURL(tc.inputURL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Test %v - '%s' FAIL: expected error, got URL: %v", i, tc.name, actual)
				}
				return
			}
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
