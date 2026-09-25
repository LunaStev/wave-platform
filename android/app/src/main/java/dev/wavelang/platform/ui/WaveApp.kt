package dev.wavelang.platform.ui

import androidx.activity.compose.BackHandler
import androidx.annotation.StringRes
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.selection.selectable
import androidx.compose.foundation.selection.selectableGroup
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.List
import androidx.compose.material.icons.filled.Info
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.saveable.rememberSaveableStateHolder
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.heading
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import dev.wavelang.platform.BuildConfig
import dev.wavelang.platform.R
import dev.wavelang.platform.data.DocumentSummary
import dev.wavelang.platform.data.DocumentationProject
import dev.wavelang.platform.data.documentLocales

private enum class Destination(@param:StringRes val label: Int, val icon: ImageVector) {
    Docs(R.string.nav_docs, Icons.AutoMirrored.Filled.List),
    News(R.string.nav_news, Icons.Default.Info),
    Settings(R.string.nav_settings, Icons.Default.Settings),
}
private val languageNames = mapOf(
    "en" to "English", "ko" to "한국어", "ja" to "日本語", "zh" to "简体中文",
    "es" to "Español", "de" to "Deutsch", "ru" to "Русский",
    "id" to "Bahasa Indonesia · Melayu", "vi" to "Tiếng Việt",
)

@Composable
fun WaveApp(
    state: CatalogState,
    onProject: (DocumentationProject) -> Unit,
    onLocale: (String) -> Unit,
    onRetry: () -> Unit,
    theme: ThemePreference,
    onTheme: (ThemePreference) -> Unit,
    onDocument: (DocumentSummary) -> Unit,
    onBlog: () -> Unit,
) {
    var destination by rememberSaveable { mutableStateOf(Destination.Docs) }
    var languagePicker by rememberSaveable { mutableStateOf(false) }
    val screenState = rememberSaveableStateHolder()
    BackHandler(enabled = destination != Destination.Docs) { destination = Destination.Docs }
    Surface(Modifier.fillMaxSize()) {
        BoxWithConstraints(Modifier.fillMaxSize().windowInsetsPadding(WindowInsets.safeDrawing)) {
            val expanded = maxWidth >= 600.dp
            Column {
                Row(Modifier.fillMaxWidth().heightIn(min = 56.dp).padding(horizontal = 16.dp, vertical = 10.dp),
                    verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(9.dp)) {
                    Image(painterResource(R.drawable.wave_logo), contentDescription = null, modifier = Modifier.size(32.dp))
                    Text(stringResource(R.string.app_name), style = MaterialTheme.typography.titleLarge)
                    VerticalDivider(Modifier.height(18.dp))
                    Text(stringResource(destination.label), style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
                HorizontalDivider()
                Row(Modifier.weight(1f)) {
                    if (expanded) {
                        Column(Modifier.width(96.dp).fillMaxHeight().padding(top = 8.dp).selectableGroup()) {
                            Destination.entries.forEach { item ->
                                DestinationItem(item, destination == item, { destination = item }, Modifier.fillMaxWidth(), true)
                            }
                        }
                        VerticalDivider()
                    }
                    Box(Modifier.weight(1f).fillMaxHeight(), contentAlignment = Alignment.TopCenter) {
                        screenState.SaveableStateProvider(destination.name) {
                            when (destination) {
                                Destination.Docs -> CatalogScreen(state, onProject, { languagePicker = true }, onRetry, onDocument)
                                Destination.News -> NewsScreen(onBlog)
                                Destination.Settings -> SettingsScreen(state.locale, { languagePicker = true }, theme, onTheme)
                            }
                        }
                    }
                }
                if (!expanded) {
                    HorizontalDivider()
                    Row(Modifier.fillMaxWidth().selectableGroup()) {
                        Destination.entries.forEach { item ->
                            DestinationItem(item, destination == item, { destination = item }, Modifier.weight(1f), false)
                        }
                    }
                }
            }
        }
    }
    if (languagePicker) AlertDialog(
        onDismissRequest = { languagePicker = false },
        containerColor = MaterialTheme.colorScheme.surface,
        title = { Text(stringResource(R.string.choose_language), style = MaterialTheme.typography.titleLarge) },
        text = {
            LazyColumn(Modifier.selectableGroup()) {
                items(documentLocales) { code ->
                    Row(Modifier.fillMaxWidth().heightIn(min = 48.dp).selectable(
                        selected = code == state.locale, role = Role.RadioButton,
                        onClick = { onLocale(code); languagePicker = false },
                    ).padding(vertical = 8.dp), verticalAlignment = Alignment.CenterVertically) {
                        RadioButton(selected = code == state.locale, onClick = null)
                        Text(languageNames.getValue(code), Modifier.padding(start = 12.dp))
                    }
                }
            }
        },
        confirmButton = { TextButton(onClick = { languagePicker = false }) { Text(stringResource(R.string.dismiss)) } },
    )
}

@Composable
private fun DestinationItem(item: Destination, selected: Boolean, onClick: () -> Unit, modifier: Modifier, rail: Boolean) {
    Box(modifier.testTag("nav-${item.name}").selectable(selected, role = Role.Tab, onClick = onClick)
        .background(if (selected && rail) MaterialTheme.colorScheme.primaryContainer else Color.Transparent)) {
        Column(Modifier.fillMaxWidth().heightIn(min = 64.dp).padding(horizontal = 8.dp, vertical = 10.dp),
            horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.spacedBy(4.dp)) {
            Icon(item.icon, contentDescription = null, modifier = Modifier.size(20.dp),
                tint = if (selected) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant)
            Text(stringResource(item.label), style = MaterialTheme.typography.labelMedium,
                fontWeight = if (selected) FontWeight.Bold else FontWeight.Normal)
        }
        if (selected) Box(Modifier.matchParentSize()) {
            Box((if (rail) Modifier.fillMaxHeight().width(3.dp) else Modifier.fillMaxWidth().height(3.dp))
                .background(MaterialTheme.colorScheme.secondary))
        }
    }
}

@Composable
private fun LanguageButton(locale: String, onLanguage: () -> Unit, modifier: Modifier = Modifier) {
    OutlinedButton(onClick = onLanguage, modifier = modifier.heightIn(min = 48.dp), shape = MaterialTheme.shapes.small,
        colors = ButtonDefaults.outlinedButtonColors(contentColor = MaterialTheme.colorScheme.onSurface),
        contentPadding = PaddingValues(horizontal = 12.dp, vertical = 10.dp)) {
        Text(languageNames.getValue(locale), Modifier.weight(1f), style = MaterialTheme.typography.bodyLarge)
        Text("▾", Modifier.padding(start = 12.dp))
    }
}

@Composable
private fun CatalogScreen(state: CatalogState, onProject: (DocumentationProject) -> Unit, onLanguage: () -> Unit,
    onRetry: () -> Unit, onDocument: (DocumentSummary) -> Unit) {
    val documents = state.visibleDocuments
    LazyColumn(Modifier.widthIn(max = 960.dp).fillMaxSize()) {
        item {
            Column(Modifier.fillMaxWidth().background(MaterialTheme.colorScheme.surfaceVariant)) {
                Column(Modifier.padding(horizontal = 16.dp, vertical = 12.dp), verticalArrangement = Arrangement.spacedBy(6.dp)) {
                    Text(stringResource(R.string.reading_language), style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant)
                    LanguageButton(state.locale, onLanguage, Modifier.fillMaxWidth())
                }
                Row(Modifier.fillMaxWidth().padding(horizontal = 16.dp).selectableGroup()) {
                    DocumentationProject.entries.forEach { project ->
                        val selected = state.project == project
                        Column(Modifier.weight(1f).selectable(selected, role = Role.Tab, onClick = { onProject(project) }),
                            horizontalAlignment = Alignment.CenterHorizontally) {
                            Box(Modifier.heightIn(min = 48.dp).padding(horizontal = 20.dp, vertical = 12.dp), contentAlignment = Alignment.Center) {
                                Text(project.name, style = MaterialTheme.typography.bodyLarge,
                                    fontWeight = if (selected) FontWeight.Bold else FontWeight.Normal,
                                    color = if (selected) MaterialTheme.colorScheme.onSurface else MaterialTheme.colorScheme.onSurfaceVariant)
                            }
                            Box(Modifier.fillMaxWidth().height(3.dp)
                                .background(if (selected) MaterialTheme.colorScheme.secondary else Color.Transparent))
                        }
                    }
                }
                HorizontalDivider()
            }
        }
        item {
            Column(Modifier.padding(horizontal = 16.dp).padding(top = 24.dp, bottom = 18.dp),
                verticalArrangement = Arrangement.spacedBy(6.dp)) {
                Text(stringResource(R.string.docs_title, state.project.name), style = MaterialTheme.typography.headlineMedium,
                    modifier = Modifier.semantics { heading() })
                Text(stringResource(if (state.project == DocumentationProject.Wave) R.string.docs_lead else R.string.whale_lead),
                    style = MaterialTheme.typography.bodyLarge, color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
            HorizontalDivider(Modifier.padding(horizontal = 16.dp), color = MaterialTheme.colorScheme.outline)
        }
        when {
            state.loading -> item {
                Row(Modifier.fillMaxWidth().padding(24.dp), horizontalArrangement = Arrangement.spacedBy(12.dp),
                    verticalAlignment = Alignment.CenterVertically) {
                    CircularProgressIndicator(Modifier.size(20.dp), strokeWidth = 2.dp)
                    Text(stringResource(R.string.loading), style = MaterialTheme.typography.bodyLarge)
                }
            }
            state.failed -> item { StatusPanel(R.string.load_error_title, R.string.load_error_body, R.string.retry, onRetry) }
            documents.isEmpty() -> item { StatusPanel(R.string.empty_title, R.string.empty_body) }
            else -> {
                item { Text(stringResource(R.string.catalog_preview), Modifier.padding(horizontal = 16.dp, vertical = 12.dp),
                    style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                documents.groupBy { it.group }.forEach { (group, entries) ->
                    item(key = "group:$group") {
                        Text(groupLabel(group), Modifier.padding(horizontal = 16.dp).padding(top = 20.dp, bottom = 8.dp)
                            .semantics { heading() }, style = MaterialTheme.typography.titleMedium)
                    }
                    items(entries, key = { it.path }) { document ->
                        Column(Modifier.fillMaxWidth().clickable(onClickLabel = stringResource(R.string.open_on_web)) {
                            onDocument(document)
                        }.padding(horizontal = 16.dp)) {
                            Column(Modifier.padding(horizontal = 8.dp, vertical = 12.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                                Text(document.title, style = MaterialTheme.typography.titleSmall)
                                if (document.summary.isNotBlank()) Text(document.summary, style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.onSurfaceVariant)
                                if (document.locale != state.locale) Text(stringResource(R.string.english_fallback),
                                    style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.primary)
                                Text(stringResource(R.string.open_on_web), style = MaterialTheme.typography.bodySmall,
                                    color = MaterialTheme.colorScheme.primary)
                            }
                            HorizontalDivider()
                        }
                    }
                }
                item { Spacer(Modifier.height(24.dp)) }
            }
        }
    }
}

@Composable
private fun groupLabel(group: String): String = when (group) {
    "getting-started" -> stringResource(R.string.group_getting_started)
    "language" -> stringResource(R.string.group_language)
    "reference" -> stringResource(R.string.group_reference)
    "toolchain" -> stringResource(R.string.group_tools)
    "whale" -> "Whale"
    "" -> stringResource(R.string.catalog_label)
    else -> group
}

@Composable
private fun StatusPanel(@StringRes title: Int, @StringRes body: Int, @StringRes action: Int? = null, onAction: () -> Unit = {}) {
    Column(Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 24.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
        Text(stringResource(title), style = MaterialTheme.typography.titleMedium)
        Text(stringResource(body), style = MaterialTheme.typography.bodyLarge, color = MaterialTheme.colorScheme.onSurfaceVariant)
        if (action != null) OutlinedButton(onClick = onAction, shape = MaterialTheme.shapes.small,
            modifier = Modifier.heightIn(min = 48.dp)) { Text(stringResource(action)) }
    }
}

@Composable
private fun SectionHeading(@StringRes title: Int) {
    Text(stringResource(title), Modifier.fillMaxWidth().padding(vertical = 18.dp).semantics { heading() },
        style = MaterialTheme.typography.headlineMedium)
    HorizontalDivider(color = MaterialTheme.colorScheme.outline)
}

@Composable
private fun NewsScreen(onBlog: () -> Unit) {
    LazyColumn(Modifier.widthIn(max = 960.dp).fillMaxSize(), contentPadding = PaddingValues(16.dp)) {
        item { SectionHeading(R.string.news_title) }
        item { StatusPanel(R.string.news_empty_title, R.string.news_empty_body, R.string.open_blog, onBlog) }
    }
}

@Composable
private fun SettingsScreen(locale: String, onLanguage: () -> Unit, theme: ThemePreference, onTheme: (ThemePreference) -> Unit) {
    LazyColumn(Modifier.widthIn(max = 960.dp).fillMaxSize(), contentPadding = PaddingValues(16.dp)) {
        item { SectionHeading(R.string.settings_title) }
        item {
            Column(Modifier.padding(vertical = 20.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(stringResource(R.string.reading_language), style = MaterialTheme.typography.titleMedium)
                Text(stringResource(R.string.reading_language_hint), style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
                LanguageButton(locale, onLanguage, Modifier.fillMaxWidth())
            }
            HorizontalDivider()
        }
        item {
            Column(Modifier.selectableGroup().padding(vertical = 20.dp)) {
                Text(stringResource(R.string.appearance), style = MaterialTheme.typography.titleMedium)
                ThemePreference.entries.forEach { preference ->
                    val label = when (preference) {
                        ThemePreference.System -> R.string.theme_system
                        ThemePreference.Light -> R.string.theme_light
                        ThemePreference.Dark -> R.string.theme_dark
                    }
                    Row(Modifier.fillMaxWidth().heightIn(min = 48.dp).selectable(theme == preference, role = Role.RadioButton,
                        onClick = { onTheme(preference) }).padding(vertical = 8.dp), verticalAlignment = Alignment.CenterVertically) {
                        RadioButton(theme == preference, onClick = null)
                        Text(stringResource(label), Modifier.padding(start = 12.dp), style = MaterialTheme.typography.bodyLarge)
                    }
                }
            }
            HorizontalDivider()
        }
        item {
            Column(Modifier.padding(vertical = 20.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(stringResource(R.string.about), style = MaterialTheme.typography.titleMedium)
                Text(stringResource(R.string.about_body), style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
                Text(stringResource(R.string.version, BuildConfig.VERSION_NAME), style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        }
    }
}
