export async function fetchWithCache(url, options = {}) {
  const cacheKey = `cache:${url}`;
  const cached = localStorage.getItem(cacheKey);

  if (cached) {
    const { etag, data, expires } = JSON.parse(cached);
    if (expires > Date.now()) {
      return data;
    }
    options.headers = {
      ...options.headers,
      "If-None-Match": etag,
    };
  }

  const response = await fetch(url, options);

  if (response.status === 304) {
    return JSON.parse(cached).data;
  }

  if (response.ok) {
    const data = await response.json();
    const etag = response.headers.get("ETag");
    const cacheControl = response.headers.get("Cache-Control");

    if (etag && cacheControl) {
      const maxAge = cacheControl.match(/max-age=(\d+)/)?.[1] || 0;
      localStorage.setItem(
        cacheKey,
        JSON.stringify({
          etag,
          data,
          expires: Date.now() + maxAge * 1000,
        }),
      );
    }
    return data;
  }

  throw new Error(`Request failed: ${response.status}`);
}
