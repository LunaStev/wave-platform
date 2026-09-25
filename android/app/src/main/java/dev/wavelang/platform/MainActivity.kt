package dev.wavelang.platform

import android.content.ActivityNotFoundException
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.SideEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.graphics.luminance
import androidx.core.view.WindowCompat
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.createSavedStateHandle
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import dev.wavelang.platform.data.HttpCatalogRepository
import dev.wavelang.platform.data.initialDocumentLocale
import dev.wavelang.platform.ui.CatalogViewModel
import dev.wavelang.platform.ui.ThemePreference
import dev.wavelang.platform.ui.WaveApp
import dev.wavelang.platform.ui.WaveTheme
import java.util.Locale

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        val preferences = getSharedPreferences("reading", MODE_PRIVATE)
        val initialLocale = initialDocumentLocale(preferences.getString("locale", null), Locale.getDefault().language)
        val factory = viewModelFactory {
            initializer {
                CatalogViewModel(HttpCatalogRepository(BuildConfig.API_BASE_URL), createSavedStateHandle(), initialLocale) {
                    preferences.edit().putString("locale", it).apply()
                }
            }
        }
        setContent {
            val model: CatalogViewModel = viewModel(factory = factory)
            val state by model.state.collectAsStateWithLifecycle()
            var theme by remember {
                mutableStateOf(ThemePreference.entries.firstOrNull { it.name == preferences.getString("theme", null) }
                    ?: ThemePreference.System)
            }
            WaveTheme(theme) {
                val lightSurface = MaterialTheme.colorScheme.background.luminance() > 0.5f
                SideEffect {
                    WindowCompat.getInsetsController(window, window.decorView).apply {
                        isAppearanceLightStatusBars = lightSurface
                        isAppearanceLightNavigationBars = lightSurface
                    }
                }
                WaveApp(state, model::selectProject, model::selectLocale, model::reload, theme,
                    onTheme = { theme = it; preferences.edit().putString("theme", it.name).apply() },
                    onDocument = { document ->
                        val uri = Uri.parse(BuildConfig.API_BASE_URL).buildUpon()
                            .appendPath("docs").appendPath(document.locale)
                        document.path.split('/').forEach(uri::appendPath)
                        openBrowser(uri.build())
                    },
                    onBlog = { openBrowser(Uri.parse(BuildConfig.API_BASE_URL).buildUpon().appendPath("blog").build()) },
                )
            }
        }
    }

    private fun openBrowser(uri: Uri) {
        try {
            startActivity(Intent(Intent.ACTION_VIEW, uri))
        } catch (_: ActivityNotFoundException) {
            Toast.makeText(this, R.string.browser_unavailable, Toast.LENGTH_LONG).show()
        }
    }
}
