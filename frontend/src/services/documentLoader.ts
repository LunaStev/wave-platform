import type { DocumentLocale, DocumentSummary, DocumentView } from './http.ts'

type Transport = {
  list: (locale: DocumentLocale) => Promise<DocumentSummary[]>
  get: (path: string, locale: DocumentLocale) => Promise<DocumentView>
  isNotFound: (error: unknown) => boolean
}

// Only successful public catalogues are cached. Locale keys are bounded by the
// supported registry; errors never prevent the next request from retrying.
export function createDocumentLoader(transport: Transport, now = Date.now) {
  const catalogs = new Map<DocumentLocale, {
    value?: DocumentSummary[]
    expires: number
    pending?: Promise<DocumentSummary[]>
  }>()
  function catalog(locale: DocumentLocale, fresh: boolean) {
    const previous = catalogs.get(locale)
    if (previous?.pending) return previous.pending
    if (!fresh && previous?.value && previous.expires > now()) return Promise.resolve(previous.value)
    const entry: { value?: DocumentSummary[]; expires: number; pending?: Promise<DocumentSummary[]> } = { expires: 0 }
    entry.pending = Promise.resolve().then(() => transport.list(locale)).then(items => {
      entry.value = items
      entry.expires = now() + 60_000
      return items
    }).catch(error => {
      catalogs.delete(locale)
      throw error
    }).finally(() => { entry.pending = undefined })
    catalogs.set(locale, entry)
    return entry.pending
  }
  async function detail(path: string, locale: DocumentLocale) {
    if (!path) return { document: null, englishFallback: false }
    try {
      return { document: await transport.get(path, locale), englishFallback: false }
    } catch (error) {
      if (locale === 'en' || !transport.isNotFound(error)) throw error
      return { document: await transport.get(path, 'en'), englishFallback: true }
    }
  }
  return {
    async load(path: string, locale: DocumentLocale, fresh = false) {
      const [translated, english, result] = await Promise.all([
        catalog(locale, fresh), locale === 'en' ? Promise.resolve([]) : catalog('en', fresh), detail(path, locale),
      ])
      return { ...result, translated, english }
    },
  }
}
