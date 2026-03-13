import { ref } from "vue";
import { requestInfo, requestScreenshotLinks, requestScreenshotZip } from "../api/media";
import { copyText, saveBlob } from "../utils/output";

export function useMediaActions(path, screenshotVariant, hasInput) {
    const output = ref("就绪。");
    const busy = ref(false);
    const copyLabel = ref("复制");

    const setBusy = (isBusy, label) => {
        busy.value = isBusy;
        if (label) {
            output.value = label;
        }
    };

    const appendOutput = (text) => {
        output.value = text;
    };

    const errorOutput = (message) => {
        output.value = `错误：${message}`;
    };

    const runInfo = async (url, label) => {
        if (!hasInput.value) {
            errorOutput("请先选择媒体路径。");
            return;
        }
        try {
            setBusy(true, `${label} 生成中...`);
            const data = await requestInfo(path.value.trim(), url);
            appendOutput(data.output || "没有输出。");
        } catch (err) {
            errorOutput(err?.message || "请求失败。");
        } finally {
            setBusy(false);
        }
    };

    const downloadShots = async () => {
        if (!hasInput.value) {
            errorOutput("请先选择媒体路径。");
            return;
        }
        try {
            setBusy(true, "正在生成截图...");
            const blob = await requestScreenshotZip(path.value.trim(), screenshotVariant.value);
            saveBlob(blob, "screenshots.zip");
            appendOutput("截图已下载为 screenshots.zip。");
        } catch (err) {
            errorOutput(err?.message || "截图请求失败。");
        } finally {
            setBusy(false);
        }
    };

    const outputShotLinks = async () => {
        if (!hasInput.value) {
            errorOutput("请先选择媒体路径。");
            return;
        }
        try {
            setBusy(true, "正在生成截图并上传...");
            const data = await requestScreenshotLinks(path.value.trim(), screenshotVariant.value);
            appendOutput(data.output || "没有返回图床链接。");
        } catch (err) {
            errorOutput(err?.message || "图床链接请求失败。");
        } finally {
            setBusy(false);
        }
    };

    const clearOutput = () => {
        if (busy.value) {
            return;
        }
        appendOutput("就绪。");
    };

    const copyOutput = async () => {
        const text = output.value || "";
        if (text.trim() === "") {
            errorOutput("没有可复制的内容。");
            return;
        }

        await copyText(text);
        const original = copyLabel.value;
        copyLabel.value = "已复制";
        setTimeout(() => {
            copyLabel.value = original;
        }, 1200);
    };

    return {
        output,
        busy,
        copyLabel,
        runInfo,
        downloadShots,
        outputShotLinks,
        clearOutput,
        copyOutput,
    };
}
