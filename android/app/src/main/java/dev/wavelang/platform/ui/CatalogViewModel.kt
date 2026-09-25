package dev.wavelang.platform.ui

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dev.wavelang.platform.data.CatalogRepository
import dev.wavelang.platform.data.DocumentSummary
import dev.wavelang.platform.data.DocumentationProject
import dev.wavelang.platform.data.documentLocales
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.Job
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class CatalogState(
    val project: DocumentationProject = DocumentationProject.Wave,
    val locale: String = "en",
    val documents: List<DocumentSummary> = emptyList(),
    val loading: Boolean = true,
    val failed: Boolean = false,
) {
    val visibleDocuments: List<DocumentSummary> get() = documents.filter { it.project == project }
}

class CatalogViewModel(
    private val repository: CatalogRepository,
    private val savedState: SavedStateHandle,
    initialLocale: String,
    private val persistLocale: (String) -> Unit = {},
) : ViewModel() {
    private val mutableState = MutableStateFlow(CatalogState(
        locale = savedState.get<String>("locale")?.takeIf { it in documentLocales } ?: initialLocale,
        project = DocumentationProject.entries.firstOrNull { it.name == savedState.get<String>("project") }
            ?: DocumentationProject.Wave,
    ))
    val state = mutableState.asStateFlow()
    private var request: Job? = null
    private var generation = 0

    init { reload() }

    fun selectProject(project: DocumentationProject) {
        savedState["project"] = project.name
        mutableState.update { it.copy(project = project) }
    }

    fun selectLocale(locale: String) {
        if (locale !in documentLocales || locale == state.value.locale) return
        savedState["locale"] = locale
        persistLocale(locale)
        mutableState.update { it.copy(locale = locale) }
        reload()
    }

    fun reload() {
        val current = ++generation
        val locale = state.value.locale
        request?.cancel()
        mutableState.update { it.copy(loading = true, failed = false, documents = emptyList()) }
        request = viewModelScope.launch {
            try {
                val documents = repository.load(locale)
                if (current == generation) mutableState.update { it.copy(documents = documents, loading = false) }
            } catch (cancelled: CancellationException) {
                throw cancelled
            } catch (_: Exception) {
                if (current == generation) mutableState.update { it.copy(loading = false, failed = true) }
            }
        }
    }
}
