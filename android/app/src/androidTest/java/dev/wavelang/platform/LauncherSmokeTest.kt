package dev.wavelang.platform

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.compose.ui.test.onNodeWithTag
import androidx.compose.ui.test.onNodeWithText
import org.junit.Rule
import org.junit.Test

/** Exercise the actual launcher and dependency wiring, even with the local API unavailable. */
class LauncherSmokeTest {
    @get:Rule val compose = createAndroidComposeRule<MainActivity>()

    @Test fun launcherStartsTheNativeShell() {
        compose.onNodeWithTag("nav-Docs").assertIsDisplayed()
        compose.onNodeWithText("Wave Documentation").assertIsDisplayed()
    }
}
