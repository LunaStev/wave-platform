import localeRegistry from '../../../wavedoc/locales.json' with { type: 'json' }
import type { DocumentLocale } from './http'

export const documentLocales: Array<{ id: DocumentLocale; label: string }> = localeRegistry.map(
  ({ id, label }) => ({ id: id as DocumentLocale, label }),
)

const supported = new Set<DocumentLocale>(documentLocales.map((item) => item.id))

export function isDocumentLocale(value: unknown): value is DocumentLocale {
  return typeof value === 'string' && supported.has(value as DocumentLocale)
}

export function initialDocumentLocale(): DocumentLocale {
  const saved = localStorage.getItem('wave-doc-locale')
  if (isDocumentLocale(saved)) return saved
  const browser = navigator.language.toLowerCase().split('-')[0]
  if (browser === 'ms') return 'id'
  return isDocumentLocale(browser) ? browser : 'en'
}

export function saveDocumentLocale(locale: DocumentLocale) {
  localStorage.setItem('wave-doc-locale', locale)
}
