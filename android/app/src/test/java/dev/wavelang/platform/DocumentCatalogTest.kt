package dev.wavelang.platform

import dev.wavelang.platform.data.*
import org.junit.Assert.*
import org.junit.Test

class DocumentCatalogTest {
    private fun doc(path: String, locale: String = "en", group: Int = 0, order: Int = 0) =
        DocumentSummary(path, locale, "A guide", "", "Tools", group, order)

    @Test fun translationsReplaceEnglishBeforeOrderingAndProjectSelection() {
        val english = listOf(doc("whale/a", group = 1), doc("whale/b", group = 2), doc("toolchain/whale-build"))
        val translated = listOf(doc("whale/a", "ko", group = 3))
        val result = mergeCatalogs(english, translated)
        assertEquals(listOf("toolchain/whale-build", "whale/b", "whale/a"), result.map { it.path })
        assertEquals("ko", result.last().locale)
        assertEquals(DocumentationProject.Wave, result.first().project)
        assertEquals(2, result.count { it.project == DocumentationProject.Whale })
    }

    @Test fun namespaceAndOptionalUnknownFieldsAreHandled() {
        val result = parseCatalog("""
            <p:documents xmlns:p="$API_NAMESPACE"><p:document>
              <p:path>whale/overview</p:path><p:locale>en</p:locale><p:title>Whale &amp; tools</p:title>
              <p:summary><![CDATA[Build <tools>]]></p:summary><p:group-order>4</p:group-order>
              <p:future><p:path>not-a-document</p:path></p:future>
            </p:document></p:documents>
        """.trimIndent())
        assertEquals(1, result.size)
        assertEquals("Whale & tools", result.single().title)
        assertEquals("Build <tools>", result.single().summary)
        assertEquals(4, result.single().groupOrder)
        assertTrue(parseCatalog("<documents xmlns=\"$API_NAMESPACE\"/>").isEmpty())
    }

    @Test fun malformedOrUntrustedXmlIsNotAnEmptySuccessfulCatalog() {
        for (xml in listOf(
            "<html>Proxy error</html>",
            "<documents xmlns=\"urn:wrong\"/>",
            "<documents xmlns=\"$API_NAMESPACE\"><document><path>../private</path><locale>en</locale><title>Bad</title></document></documents>",
            "<!DOCTYPE documents [<!ENTITY x SYSTEM 'file:///nonexistent'>]><documents xmlns=\"$API_NAMESPACE\">&x;</documents>",
            "<documents xmlns=\"$API_NAMESPACE\"><document>",
        )) {
            assertThrows(xml, Exception::class.java) { parseCatalog(xml) }
        }
    }

    @Test fun languagePreferenceSupportsMalayAndRejectsUnknownSavedValues() {
        assertEquals("ko", initialDocumentLocale("ko", "en"))
        assertEquals("id", initialDocumentLocale(null, "ms"))
        assertEquals("ja", initialDocumentLocale("invalid", "ja"))
        assertEquals("en", initialDocumentLocale(null, "fr"))
    }
}
