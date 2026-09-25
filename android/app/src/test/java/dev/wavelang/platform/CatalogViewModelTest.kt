package dev.wavelang.platform

import androidx.lifecycle.SavedStateHandle
import dev.wavelang.platform.data.*
import dev.wavelang.platform.ui.CatalogViewModel
import java.io.IOException
import kotlinx.coroutines.*
import kotlinx.coroutines.test.*
import org.junit.After
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class CatalogViewModelTest {
    private val dispatcher = StandardTestDispatcher()
    @Before fun setMain() { Dispatchers.setMain(dispatcher) }
    @After fun resetMain() { Dispatchers.resetMain() }

    private fun doc(locale: String) = DocumentSummary("whale/overview", locale, locale, "", "", 0, 0)

    @Test fun lateRequestsCannotReplaceANewerLanguageOrProject() = runTest(dispatcher) {
        val english = CompletableDeferred<List<DocumentSummary>>()
        val korean = CompletableDeferred<List<DocumentSummary>>()
        val repository = CatalogRepository { locale ->
            // Deliberately emulate a transport that returns after cancellation.
            withContext(NonCancellable) { if (locale == "en") english.await() else korean.await() }
        }
        val saved = SavedStateHandle()
        val model = CatalogViewModel(repository, saved, "en")
        runCurrent()
        model.selectLocale("ko")
        model.selectProject(DocumentationProject.Whale)
        runCurrent()
        korean.complete(listOf(doc("ko")))
        runCurrent()
        english.complete(listOf(doc("en")))
        runCurrent()
        assertEquals("ko", model.state.value.locale)
        assertEquals("ko", model.state.value.visibleDocuments.single().locale)
        assertEquals(DocumentationProject.Whale, model.state.value.project)
        assertEquals("ko", saved.get<String>("locale"))
        assertEquals("Whale", saved.get<String>("project"))
    }

    @Test fun failedLoadIsRetryableAndEmptyIsASuccess() = runTest(dispatcher) {
        var calls = 0
        val model = CatalogViewModel(CatalogRepository {
            if (calls++ == 0) throw IOException("offline") else emptyList()
        }, SavedStateHandle(), "en")
        runCurrent()
        assertTrue(model.state.value.failed)
        assertFalse(model.state.value.loading)
        model.reload()
        runCurrent()
        assertFalse(model.state.value.failed)
        assertTrue(model.state.value.documents.isEmpty())
        assertEquals(2, calls)
    }

    @Test fun recreationUsesSavedProjectAndLanguage() = runTest(dispatcher) {
        val saved = SavedStateHandle(mapOf("locale" to "ja", "project" to "Whale"))
        val model = CatalogViewModel(CatalogRepository { emptyList() }, saved, "en")
        runCurrent()
        assertEquals("ja", model.state.value.locale)
        assertEquals(DocumentationProject.Whale, model.state.value.project)
    }
}
