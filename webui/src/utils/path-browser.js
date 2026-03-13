export function normalizeComparePath(value) {
    if (!value) {
        return "";
    }
    if (value === "/" || value === "\\") {
        return "/";
    }
    return value.replace(/\\/g, "/").replace(/\/+$/, "").toLowerCase();
}

export function withTrailingSeparator(value) {
    if (value === "") {
        return "";
    }
    if (value.endsWith("/") || value.endsWith("\\")) {
        return value;
    }
    const separator = value.includes("\\") && !value.includes("/") ? "\\" : "/";
    return `${value}${separator}`;
}

export function cleanPath(value) {
    if (!value) {
        return "";
    }
    if (value === "/" || value === "\\") {
        return value;
    }
    return value.replace(/[\\/]+$/, "");
}

export function getEntryName(value) {
    const normalized = value.replace(/[\\/]+$/, "");
    if (normalized === "") {
        return value;
    }
    const parts = normalized.split(/[\\/]/);
    return parts[parts.length - 1] || normalized;
}

export function buildEntries(items) {
    const result = [];
    for (const raw of items) {
        if (typeof raw !== "string" || raw.trim() === "") {
            continue;
        }
        const isDir = raw.endsWith("/") || raw.endsWith("\\");
        result.push({
            path: cleanPath(raw),
            name: getEntryName(raw),
            isDir,
        });
    }

    result.sort((left, right) => {
        if (left.isDir !== right.isDir) {
            return left.isDir ? -1 : 1;
        }
        return left.name.localeCompare(right.name, "zh-CN");
    });

    return result;
}

export function filterEntries(entries, keyword) {
    const normalizedKeyword = keyword.trim().toLowerCase();
    if (normalizedKeyword === "") {
        return entries;
    }

    return entries.filter((entry) => {
        const name = (entry.name || "").toLowerCase();
        const full = (entry.path || "").toLowerCase();
        return name.includes(normalizedKeyword) || full.includes(normalizedKeyword);
    });
}

export function getParentDirectory(dir, root) {
    const normalized = cleanPath(dir);
    if (normalized === "" || normalized === "/") {
        return normalized;
    }

    const slash = Math.max(normalized.lastIndexOf("/"), normalized.lastIndexOf("\\"));
    if (slash <= 0) {
        return root || "";
    }
    return normalized.slice(0, slash);
}

export function canNavigateUp(browserDir, browserRoot, browserRoots) {
    if (!browserDir) {
        return false;
    }

    const root = normalizeComparePath(browserRoot);
    const current = normalizeComparePath(browserDir);
    if (root === "") {
        return true;
    }
    if (current !== root) {
        return true;
    }
    return browserRoots.length > 1;
}
