package dev.wavelang.platform.ui

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

// Match frontend/src/ui/tokens.css. Android keeps its platform sans-serif fallback.
private val LightColors = lightColorScheme(
    primary = Color(0xFF4937C7), onPrimary = Color.White,
    primaryContainer = Color(0xFFF1EFFF), onPrimaryContainer = Color(0xFF4938BA),
    secondary = Color(0xFF6654F1), onSecondary = Color.White,
    background = Color(0xFFFFFFFF), onBackground = Color(0xFF1F2328),
    surface = Color(0xFFFFFFFF), onSurface = Color(0xFF1F2328),
    surfaceVariant = Color(0xFFF8F9FB), onSurfaceVariant = Color(0xFF5F6672),
    outline = Color(0xFFB9BEC9), outlineVariant = Color(0xFFD8DBE2),
    surfaceTint = Color.Transparent,
)
private val DarkColors = darkColorScheme(
    primary = Color(0xFFA99FFF), onPrimary = Color(0xFF15171C),
    primaryContainer = Color(0xFF24203F), onPrimaryContainer = Color(0xFFC4BDFF),
    secondary = Color(0xFF6654F1), onSecondary = Color.White,
    background = Color(0xFF15171C), onBackground = Color(0xFFF0F2F5),
    surface = Color(0xFF15171C), onSurface = Color(0xFFF0F2F5),
    surfaceVariant = Color(0xFF111318), onSurfaceVariant = Color(0xFFBAC0CA),
    outline = Color(0xFF464C58), outlineVariant = Color(0xFF2E323B),
    surfaceTint = Color.Transparent,
)
private val WaveShapes = Shapes(
    extraSmall = RoundedCornerShape(2.dp), small = RoundedCornerShape(4.dp),
    medium = RoundedCornerShape(6.dp), large = RoundedCornerShape(8.dp),
    extraLarge = RoundedCornerShape(8.dp),
)
private fun text(size: Int, line: Int, weight: FontWeight = FontWeight.Normal) =
    TextStyle(fontFamily = FontFamily.SansSerif, fontSize = size.sp, lineHeight = line.sp,
        fontWeight = weight, letterSpacing = 0.sp)
private val WaveTypography = Typography(
    headlineMedium = text(27, 34, FontWeight.Bold),
    titleLarge = text(19, 26, FontWeight.Bold),
    titleMedium = text(17, 24, FontWeight.SemiBold),
    titleSmall = text(14, 20, FontWeight.SemiBold),
    bodyLarge = text(14, 21), bodyMedium = text(13, 19), bodySmall = text(12, 18),
    labelLarge = text(14, 20, FontWeight.SemiBold),
    labelMedium = text(12, 18), labelSmall = text(11, 16),
)

enum class ThemePreference { System, Light, Dark }

@Composable
fun WaveTheme(preference: ThemePreference = ThemePreference.System, content: @Composable () -> Unit) {
    val dark = when (preference) {
        ThemePreference.System -> isSystemInDarkTheme()
        ThemePreference.Light -> false
        ThemePreference.Dark -> true
    }
    MaterialTheme(colorScheme = if (dark) DarkColors else LightColors,
        shapes = WaveShapes, typography = WaveTypography, content = content)
}
