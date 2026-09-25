package dev.wavelang.platform

import android.graphics.Bitmap
import android.os.ParcelFileDescriptor
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.graphics.asAndroidBitmap
import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.junit4.StateRestorationTester
import androidx.test.platform.app.InstrumentationRegistry
import dev.wavelang.platform.data.DocumentSummary
import dev.wavelang.platform.data.DocumentationProject
import dev.wavelang.platform.ui.*
import org.junit.Assert.assertEquals
import org.junit.Rule
import org.junit.Test
import java.io.File

class WaveAppTest {
    @get:Rule val compose = createComposeRule()

    @Test fun navigationAndSelectionSurviveRecreation() {
        val restoration = StateRestorationTester(compose)
        val state = mutableStateOf(CatalogState(loading = false))
        restoration.setContent {
            WaveTheme { WaveApp(state.value, { state.value = state.value.copy(project = it) },
                { state.value = state.value.copy(locale = it) }, {}, ThemePreference.System, {}, {}, {}) }
        }
        compose.onNodeWithText("Whale").performClick()
        compose.runOnIdle { assertEquals(DocumentationProject.Whale, state.value.project) }
        compose.onNodeWithTag("nav-Settings").performClick()
        compose.onNodeWithText("Appearance").assertIsDisplayed()
        capture("settings")
        restoration.emulateSavedInstanceStateRestore()
        compose.onNodeWithText("Appearance").assertIsDisplayed()
        compose.onNodeWithText("English").performClick()
        compose.onNodeWithText("한국어").performClick()
        compose.runOnIdle { assertEquals("ko", state.value.locale) }
    }

    @Test fun failureOffersRetryInsteadOfAnEmptyCatalog() {
        var retries = 0
        compose.setContent {
            WaveTheme { WaveApp(CatalogState(loading = false, failed = true), {}, {}, { retries++ },
                ThemePreference.System, {}, {}, {}) }
        }
        compose.onNodeWithText("Couldn’t load the guides").assertIsDisplayed()
        compose.onNodeWithText("Try again").performClick()
        compose.runOnIdle { assertEquals(1, retries) }
        compose.onNodeWithText("No guides published yet").assertDoesNotExist()
    }

    @Test fun onlyTheSelectedProjectsRealDocumentsAppear() {
        val documents = listOf(
            DocumentSummary("language/intro", "en", "Language guide", "Read Wave", "Language", 1, 1),
            DocumentSummary("whale/intro", "en", "Whale guide", "Build tools", "Whale", 1, 1),
        )
        val state = mutableStateOf(CatalogState(loading = false, documents = documents))
        compose.setContent {
            WaveTheme { WaveApp(state.value, { state.value = state.value.copy(project = it) }, {}, {},
                ThemePreference.System, {}, {}, {}) }
        }
        compose.onNodeWithText("Language guide").assertExists()
        compose.onNodeWithText("Whale guide").assertDoesNotExist()
        capture("wave-catalog")
        compose.onNodeWithText("Whale").performClick()
        compose.onNodeWithText("Whale guide").assertExists()
        compose.onNodeWithText("Language guide").assertDoesNotExist()
    }

    private fun capture(name: String) {
        val context = InstrumentationRegistry.getInstrumentation().targetContext
        val directory = File(context.getExternalFilesDir(null), "screenshots").apply { mkdirs() }
        File(directory, "$name.png").outputStream().use {
            compose.onRoot().captureToImage().asAndroidBitmap().compress(Bitmap.CompressFormat.PNG, 100, it)
        }
        // Keep review artifacts even if the test runner uninstalls the fixture app.
        val shell = InstrumentationRegistry.getInstrumentation().uiAutomation.executeShellCommand(
            "sh -c 'mkdir -p /data/local/tmp/wave-ui && cp ${directory.absolutePath}/$name.png /data/local/tmp/wave-ui/'")
        ParcelFileDescriptor.AutoCloseInputStream(shell).use { it.readBytes() }
    }
}
