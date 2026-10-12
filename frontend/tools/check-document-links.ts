import { readFileSync, readdirSync } from 'node:fs'
import { renderDocumentMarkdown } from '../src/services/documentMarkdown.ts'
import { validateDocumentLinks, type DocumentLinkEntry } from '../src/services/documentLinks.ts'

const root = new URL('../../wavedoc/', import.meta.url)
const locales = new Set<string>(JSON.parse(readFileSync(new URL('locales.json', root), 'utf8')).map(item => item.id))
const redirects = JSON.parse(readFileSync(new URL('redirects.json', root), 'utf8'))
const documents: DocumentLinkEntry[] = []
for (const name of readdirSync(root, { recursive: true }).map(String).sort()) {
  if (!name.endsWith('.md') || name.split('/').length !== 3 || !locales.has(name.split('/')[0])) continue
  const source = readFileSync(new URL(name, root), 'utf8').replace(/\r\n/g, '\n')
  const match = /^---\n([\s\S]*?)\n---\n([\s\S]*)$/.exec(source)
  if (!match) throw new Error(`${name}: missing front matter`)
  const metadata = Object.fromEntries(match[1].split('\n').map(line => {
    const colon = line.indexOf(':')
    if (colon < 1) throw new Error(`${name}: invalid front matter`)
    return [line.slice(0, colon).trim(), line.slice(colon + 1).trim()]
  }))
  if (!metadata.path || metadata.locale !== name.split('/')[0]) throw new Error(`${name}: invalid path or locale`)
  const rendered = renderDocumentMarkdown(match[2], { playgrounds: false, slotPrefix: 'check', renderCode: () => '' })
  documents.push({ locale: metadata.locale, path: metadata.path, anchors: new Set(rendered.headings.map(item => item.anchor)), links: rendered.links })
}
if (!documents.length) throw new Error('No official documents found')
const failures = validateDocumentLinks(documents, locales, redirects)
if (failures.length) { console.error(failures.join('\n')); process.exitCode = 1 }
else console.log(`Documentation paths and rendered heading links: PASS (${documents.length} documents)`)
