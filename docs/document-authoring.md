# Editing language documentation

The editable source for the official Wave documentation is stored in the root
`wavedoc/{locale}` module. Documentation supports `en`, `ko`, `ja`, `zh`, `es`,
`de`, `ru`, `id`, and `vi`. The `zh` locale is Simplified Chinese, and `id`
covers the Indonesian–Malay documentation without a separate `ms` locale.

```text
wavedoc/
├── en/language/explicit-memory-type-model.md
└── ko/language/explicit-memory-type-model.md
```

Each file contains front matter followed by Markdown:

````markdown
---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: en
group: language
group_order: 2
order: 7
title: Wave Explicit Memory Type Model
summary: Pointer types and explicit memory access in Wave.
---

## Pointer types

`ptr<T>` is dedicated syntax in the Wave Explicit Memory Type Model. It is not a general-purpose generic type.

```wave
var address: ptr<i32> = raw as ptr<i32>;
```
````

Use `##` for the first heading because the page title is rendered from the front matter. GFM tables, lists, block quotes, links, and fenced code blocks are supported.

Keep `path`, `group`, and ordering fields consistent between translations. Use the same `translation_set_id` for pages that represent the same document in different languages.

English is the explicit fallback for a path that has not been translated. Do
not copy English text into another locale merely to make the translation look
complete. Documentation describes the current compiler contract without a
manual-wide Wave version label. Local variables use `var`; `let` and `let mut`
are removed syntax.

## Wave and Whale navigation

The documentation header provides separate Wave and Whale navigation. A document
belongs to Whale when its front-matter `path` starts with `whale/`; all other
paths remain in Wave. Group names and titles do not determine the project.
Existing Wave URLs, including `toolchain/whale-*`, remain unchanged.

Add Whale Markdown files directly to `wavedoc/{locale}/whale/`:

```yaml
---
translation_set_id: whale-symbols
path: whale/symbols
locale: ko
group: whale
group_order: 5
order: 4
title: 심볼과 외부 연결 이름
summary: 함수와 전역 변수의 식별 및 외부 연결 이름 규칙입니다.
---
```

The example produces `/docs/ko/whale/symbols`. The Whale catalogue is
`/docs/ko/whale`, while the Wave catalogue stays at `/docs/ko`. Use the same path
and translation set in English at `wavedoc/en/whale/symbols.md`.

Sidebars, search, and previous/next links stay inside the selected project.
Language changes retain the project and document path. Missing translations use
the English page with the existing fallback notice; a project without documents
has an empty catalogue. Navigation is sorted by `group_order`, then `order`, then
`path`, after translations replace their corresponding English entries.

No per-document route or frontend registry is needed. Markdown is embedded in
the Go binary, so **rebuild and restart the application after adding or editing
files** (for Docker, `docker compose up --build -d`). The server imports the
embedded files at startup and stores published revisions in the platform
database. Do not edit generated XML or database values by hand.
