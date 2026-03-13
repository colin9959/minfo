<template>
    <div class="grain"></div>
    <main class="shell">
        <AppHeader />

        <section class="panel">
            <PathBrowser
                v-model:path="path"
                v-model:search-keyword="searchKeyword"
                :busy="busy"
                :browser-dir="browserDir"
                :browser-error="browserError"
                :browser-loading="browserLoading"
                :can-navigate-up="canNavigateUp"
                :entries="filteredEntries"
                @navigate-up="navigateUp"
                @refresh="refreshBrowser"
                @open-entry="handleEntryDoubleClick"
            />

            <div class="field">
                <label>截图模式</label>
                <ScreenshotVariantPicker v-model="screenshotVariant" :busy="busy" />
            </div>

            <ActionButtons
                :busy="busy"
                :has-input="hasInput"
                @mediainfo="runInfo('/api/mediainfo', 'MediaInfo')"
                @bdinfo="runInfo('/api/bdinfo', 'BDInfo')"
                @download-shots="downloadShots"
                @output-links="outputShotLinks"
            />
        </section>

        <OutputPanel
            :busy="busy"
            :copy-label="copyLabel"
            :output="output"
            @copy="copyOutput"
            @clear="clearOutput"
        />
    </main>
</template>

<script setup>
import { ref } from "vue";
import ActionButtons from "./components/ActionButtons.vue";
import AppHeader from "./components/AppHeader.vue";
import OutputPanel from "./components/OutputPanel.vue";
import PathBrowser from "./components/PathBrowser.vue";
import ScreenshotVariantPicker from "./components/ScreenshotVariantPicker.vue";
import { useMediaActions } from "./composables/useMediaActions";
import { usePathBrowser } from "./composables/usePathBrowser";

const screenshotVariant = ref("png");
const pathBrowser = usePathBrowser();
const mediaActions = useMediaActions(pathBrowser.path, screenshotVariant, pathBrowser.hasInput);

const {
    path,
    searchKeyword,
    browserDir,
    browserError,
    browserLoading,
    canNavigateUp,
    filteredEntries,
    hasInput,
    navigateUp,
    refreshBrowser,
    handleEntryDoubleClick,
} = pathBrowser;

const { output, busy, copyLabel, runInfo, downloadShots, outputShotLinks, clearOutput, copyOutput } = mediaActions;
</script>
