export type DocumentLinkEntry = { locale: string; path: string; anchors: Set<string>; links: string[] }

// Validate local documentation URLs only; no external network access is needed.
export function validateDocumentLinks(entries: DocumentLinkEntry[], locales: Set<string>, redirects: Record<string, string>) {
  const documents = new Map<string, DocumentLinkEntry>()
  const failures: string[] = []
  for (const entry of entries) {
    const key = `${entry.locale}/${entry.path}`
    if (documents.has(key)) failures.push(`${key}: duplicate document path`)
    documents.set(key, entry)
  }
  for (const entry of entries) {
    for (const href of entry.links) {
      const markdownFile = !href.startsWith('/') && /^(?![a-z]+:).*\.md(?:[#?]|$)/i.test(href)
      if (!href.startsWith('#') && !href.startsWith('/docs/') && !markdownFile) continue
      const prefix = `${entry.locale}/${entry.path}: ${href}`
      try {
        const url = new URL(href, `https://documentation.invalid/docs/${entry.locale}/${entry.path}`)
        const parts = decodeURIComponent(url.pathname).split('/').filter(Boolean)
        if (parts[0] !== 'docs') { failures.push(`${prefix}: target is outside documentation`); continue }
        const locale = parts[1]
        if (!locales.has(locale)) { failures.push(`${prefix}: unsupported locale`); continue }
        let path = parts.slice(2).join('/')
        if (/\.md$/i.test(path)) { failures.push(`${prefix}: use a /docs/ URL instead of a Markdown file link`); continue }
        path = redirects[path] ?? path
        const anchor = decodeURIComponent(url.hash.slice(1))
        if (['', 'stdlib', 'whale'].includes(path)) {
          if (anchor) failures.push(`${prefix}: catalogue anchors are not supported`)
          continue
        }
        const target = documents.get(`${locale}/${path}`) ?? documents.get(`en/${path}`)
        if (!target) failures.push(`${prefix}: document not found`)
        else if (anchor && !target.anchors.has(anchor)) failures.push(`${prefix}: heading not found`)
      } catch {
        failures.push(`${prefix}: invalid URL or fragment encoding`)
      }
    }
  }
  return failures
}
