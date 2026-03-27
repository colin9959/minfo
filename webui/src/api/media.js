export async function fetchDirectory(prefix = "", signal) {
    const url = new URL("/api/path", window.location.origin);
    if (prefix !== "") {
        url.searchParams.set("prefix", prefix);
    }

    const response = await fetch(url.toString(), { signal });
    const data = await response.json();
    if (!response.ok || !data.ok || !Array.isArray(data.items)) {
        throw new Error(data.error || "读取路径失败。");
    }
    return data;
}

export async function requestInfo(path, url, fields = {}) {
    const response = await postForm(url, { path, ...fields });
    const data = await safeReadJSON(response);
    if (!response.ok || !data.ok) {
        throw new Error(data.error || "请求失败。");
    }
    return data;
}

export async function prepareScreenshotZipDownload(path, variant) {
    const response = await postForm("/api/screenshots", { path, mode: "zip", variant, prepare_download: "1" });
    const data = await safeReadJSON(response);
    if (!response.ok || !data.ok || typeof data.output !== "string" || data.output.trim() === "") {
        throw new Error(data.error || "截图请求失败。");
    }
    return new URL(data.output, window.location.origin).toString();
}

export function startPreparedDownload(url) {
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.style.display = "none";
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
}

export async function requestScreenshotLinks(path, variant) {
    const response = await postForm("/api/screenshots", { path, mode: "links", variant });
    const data = await safeReadJSON(response);
    if (!response.ok || !data.ok) {
        throw new Error(data.error || "图床链接请求失败。");
    }
    return data;
}

async function postForm(url, fields = {}) {
    const form = new FormData();
    for (const [key, value] of Object.entries(fields)) {
        if (value !== undefined && value !== null && `${value}` !== "") {
            form.append(key, `${value}`);
        }
    }
    return fetch(url, { method: "POST", body: form });
}

async function safeReadJSON(response) {
    try {
        return await response.json();
    } catch {
        return {};
    }
}
