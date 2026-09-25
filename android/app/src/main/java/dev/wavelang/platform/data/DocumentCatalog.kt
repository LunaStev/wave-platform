package dev.wavelang.platform.data

import java.io.ByteArrayOutputStream
import java.io.IOException
import java.io.StringReader
import java.net.HttpURLConnection
import java.net.URI
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.withContext
import org.xmlpull.v1.XmlPullParser
import org.xmlpull.v1.XmlPullParserFactory

const val API_NAMESPACE = "https://wave-lang.dev/ns/platform/api/v1"
val documentLocales = listOf("en", "ko", "ja", "zh", "es", "de", "ru", "id", "vi")

enum class DocumentationProject { Wave, Whale }

data class DocumentSummary(
    val path: String,
    val locale: String,
    val title: String,
    val summary: String,
    val group: String,
    val groupOrder: Int,
    val order: Int,
) {
    val project: DocumentationProject
        get() = if (path.startsWith("whale/")) DocumentationProject.Whale else DocumentationProject.Wave
}

fun initialDocumentLocale(saved: String?, deviceLanguage: String): String =
    saved?.takeIf { it in documentLocales }
        ?: (if (deviceLanguage == "ms") "id" else deviceLanguage).takeIf { it in documentLocales }
        ?: "en"

fun mergeCatalogs(english: List<DocumentSummary>, translated: List<DocumentSummary>): List<DocumentSummary> =
    (english + translated).associateBy { it.path }.values
        .sortedWith(compareBy(DocumentSummary::groupOrder, DocumentSummary::order, DocumentSummary::path))

/** Only the anonymous catalog contract; account transport and a full reader are separate work. */
fun interface CatalogRepository {
    suspend fun load(locale: String): List<DocumentSummary>
}

class CatalogException(val status: Int? = null) : IOException("Document catalog unavailable")

class HttpCatalogRepository(private val baseUrl: String) : CatalogRepository {
    override suspend fun load(locale: String): List<DocumentSummary> = coroutineScope {
        require(locale in documentLocales)
        val english = async { fetch("en") }
        if (locale == "en") mergeCatalogs(english.await(), emptyList()) else {
            val translated = async { fetch(locale) }
            mergeCatalogs(english.await(), translated.await())
        }
    }

    private suspend fun fetch(locale: String): List<DocumentSummary> = withContext(Dispatchers.IO) {
        val url = URI(baseUrl).resolve("api/v1/documents?locale=$locale").toURL()
        val connection = url.openConnection() as HttpURLConnection
        try {
            connection.connectTimeout = 10_000
            connection.readTimeout = 15_000
            connection.instanceFollowRedirects = false
            connection.setRequestProperty("Accept", "application/xml")
            val status = connection.responseCode
            if (status != HttpURLConnection.HTTP_OK) throw CatalogException(status)
            val bytes = ByteArrayOutputStream()
            connection.inputStream.use { input ->
                val buffer = ByteArray(8192)
                while (true) {
                    val read = input.read(buffer)
                    if (read < 0) break
                    if (bytes.size() + read > 2 * 1024 * 1024) throw CatalogException()
                    bytes.write(buffer, 0, read)
                }
            }
            parseCatalog(bytes.toString(Charsets.UTF_8.name()))
        } finally {
            connection.disconnect()
        }
    }
}

fun parseCatalog(xml: String): List<DocumentSummary> {
    // Public catalogs never contain a DTD. Reject it before any parser can resolve entities.
    if (xml.contains("<!DOCTYPE")) throw CatalogException()
    val parser = XmlPullParserFactory.newInstance().apply { isNamespaceAware = true }.newPullParser()
    if (parser.getFeature(XmlPullParser.FEATURE_PROCESS_DOCDECL)) {
        parser.setFeature(XmlPullParser.FEATURE_PROCESS_DOCDECL, false)
    }
    parser.setInput(StringReader(xml))
    val documents = mutableListOf<DocumentSummary>()
    var rootSeen = false
    var rootClosed = false
    var fields: MutableMap<String, String>? = null
    var field: String? = null
    val text = StringBuilder()
    while (parser.nextToken() != XmlPullParser.END_DOCUMENT) {
        when (parser.eventType) {
            XmlPullParser.DOCDECL -> throw CatalogException()
            XmlPullParser.START_TAG -> {
                if (parser.depth == 1) {
                    if (parser.name != "documents" || parser.namespace != API_NAMESPACE) throw CatalogException()
                    rootSeen = true
                } else if (parser.depth == 2 && parser.name == "document" && parser.namespace == API_NAMESPACE) {
                    fields = mutableMapOf()
                } else if (parser.depth == 3 && fields != null && parser.namespace == API_NAMESPACE) {
                    field = parser.name
                    text.clear()
                }
            }
            XmlPullParser.TEXT, XmlPullParser.CDSECT, XmlPullParser.ENTITY_REF -> {
                if (field != null && parser.depth == 3) text.append(parser.text.orEmpty())
            }
            XmlPullParser.END_TAG -> {
                if (parser.depth == 1) {
                    rootClosed = true
                } else if (parser.depth == 3 && field != null) {
                    fields?.put(field, text.toString().trim())
                    field = null
                } else if (parser.depth == 2 && fields != null) {
                    val values = fields
                    val path = values["path"].orEmpty()
                    val locale = values["locale"].orEmpty()
                    val title = values["title"].orEmpty()
                    if (path.isBlank() || path.startsWith('/') || path.split('/').any { it == ".." || it.isEmpty() } ||
                        locale !in documentLocales || title.isBlank()) throw CatalogException()
                    documents += DocumentSummary(path, locale, title, values["summary"].orEmpty(),
                        values["group"].orEmpty(), values["group-order"]?.toIntOrNull() ?: 0,
                        values["order"]?.toIntOrNull() ?: 0)
                    fields = null
                }
            }
        }
    }
    if (!rootSeen || !rootClosed) throw CatalogException()
    return documents
}
