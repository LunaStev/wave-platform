# Search indexing

Public blog articles, release notes, and documentation are served with visible
HTML content before JavaScript runs. The Go server renders Markdown and sanitizes
the result; Vue supplies the interactive interface. Every user agent receives the
same content. Public document and article data is embedded as escaped JSON/XML to
avoid fetching it again on the initial page load. Drafts and private mailbox data
are never included.

## Canonical URLs and languages

- Releases use `/releases/{slug}`; the old blog URL permanently redirects.
- Documentation uses `/docs/{locale}/{path}` with a canonical for each published
  translation. Untranslated pages use the English canonical and content language.
- HTML and the XML sitemap link actual published translations with reciprocal
  `hreflang` entries, including an English `x-default`.
- Blog content language is detected locally from prose, excluding fenced code,
  inline code, and URLs. Low-confidence or very short text defaults to English.
  The reader's interface language does not change the article's language metadata.
  Detection is an estimate; review mixed-language and unusually short articles.
- Missing public resources return HTTP 404 and `noindex`. Storage outages return
  HTTP 503 with `Retry-After`, rather than implying that content was deleted.

## Sitemap and crawling

`/sitemap.xml` contains public canonical URLs and real modification dates where
available. All published documentation locales are included. Collections larger
than 10,000 URLs are split through a sitemap index and numbered child sitemaps.
There are no invented priorities, update dates, or crawl-frequency claims.

`robots.txt` applies to all crawlers. It permits the public read APIs required by
the frontend while excluding private API areas. API responses carry `noindex`, so
search results point to the HTML pages rather than their XML data. Authentication
and authorization are enforced independently of robots rules.

## Deployment verification

Use the deployed public origin for these checks:

1. Fetch a blog article and a translated document without JavaScript. Verify the
   actual body, canonical, `Content-Language`, Open Graph locale, and JSON-LD.
2. Enable JavaScript and navigate between documents and languages. Confirm the
   canonical and alternate links reflect the current document.
3. Fetch a nonexistent document or article and confirm HTTP 404. Fetch
   `/sitemap.xml` and verify XML content type and valid XML, not the SPA shell.
4. Validate an article with Google's Rich Results Test and Schema.org's validator.
5. Submit `/sitemap.xml` to Google Search Console and Bing Webmaster Tools for
   the verified production property. Inspect sample URLs and monitor crawling,
   indexing, selected canonicals, and field Core Web Vitals after deployment.

The platform does not register webmaster accounts or submit URLs automatically.
Local checks verify implementation; production indexing and rich-result display
are decided by each search engine.
