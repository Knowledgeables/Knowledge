import { JSDOM } from 'jsdom';
import { URL } from 'url';

const visitedUrls = new Set();
const pagesToIngest = [];

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';
const CRAWLER_KEY = process.env.CRAWLER_KEY;
const TARGETS_ENDPOINT = process.env.TARGETS_ENDPOINT || '/api/crawler/targets';
const INGEST_ENDPOINT = process.env.INGEST_ENDPOINT || '/api/crawler/ingest';

const MAX_TARGETS = Number(process.env.MAX_TARGETS || 20);
const MAX_PAGES = Number(process.env.MAX_PAGES || 50);
const MAX_DEPTH = Number(process.env.MAX_DEPTH || 1);
const CRAWL_DELAY_MS = Number(process.env.CRAWL_DELAY_MS || 1000);
const ENABLE_INGEST = String(process.env.ENABLE_INGEST || 'false').toLowerCase() === 'true';

const FALLBACK_START_URLS = (process.env.FALLBACK_START_URLS || 'https://en.wikipedia.org/wiki/Docker_(software)')
    .split(',')
    .map((url) => url.trim())
    .filter(Boolean);

function delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function unwrapWaybackUrl(urlString) {
    const waybackPattern = /^https?:\/\/web\.archive\.org\/web\/\d+\*?\/(https?:\/\/.+)$/;
    const match = urlString.match(waybackPattern);
    if (match) {
        return match[1];
    }
    return urlString;
}

function asSeedURL(rawTarget) {
    const trimmed = (rawTarget || '').trim();
    if (!trimmed) {
        return '';
    }

    if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
        return trimmed;
    }

    return `https://en.wikipedia.org/wiki/${encodeURIComponent(trimmed)}`;
}

function extractTitle(document, fallbackURL) {
    const heading = document.querySelector('h1.firstHeading') || document.querySelector('h1');
    if (heading && heading.textContent) {
        return heading.textContent.trim();
    }

    if (document.title) {
        return document.title.replace(' - Wikipedia', '').trim();
    }

    return fallbackURL;
}

function extractMainContent(document) {
    const wikiContent = Array.from(document.querySelectorAll('.mw-parser-output p, .mw-parser-output h1, .mw-parser-output h2, .mw-parser-output h3'));
    const genericContent = Array.from(document.querySelectorAll('main p, article p, h1, h2, h3'));
    const source = wikiContent.length > 0 ? wikiContent : genericContent;

    return source
        .map(element => element.textContent.trim())
        .filter(Boolean)
        .join(' ')
        .replace(/[\t\n\r]+/g, ' ')
        .trim();
}

async function getSeedURLs() {
    try {
        const url = new URL(TARGETS_ENDPOINT, BACKEND_URL);
        url.searchParams.set('limit', String(MAX_TARGETS));

        const response = await fetch(url, {
            headers: {
                'X-Crawler-Key': CRAWLER_KEY
            }
        });

        if (!response.ok) {
            throw new Error(`targets endpoint failed: ${response.status}`);
        }

        const payload = await response.json();
        const targets = Array.isArray(payload.targets) ? payload.targets : [];
        const urls = targets
            .map((target) => asSeedURL(target.query))
            .filter(Boolean);

        if (urls.length === 0) {
            console.warn('[WARN] No targets returned from API, using fallback URL list');
            return FALLBACK_START_URLS;
        }

        return urls;
    } catch (error) {
        console.warn(`[WARN] Failed to load targets from API: ${error.message}`);
        console.warn('[WARN] Falling back to static start URL list');
        return FALLBACK_START_URLS;
    }
}

async function crawlPage(pageUrl, depth = 0) {
    if (pagesToIngest.length >= MAX_PAGES) {
        return;
    }

    const cleanUrl = new URL(pageUrl);
    cleanUrl.hash = '';

    if (visitedUrls.has(cleanUrl.href)) {
        return;
    }

    visitedUrls.add(cleanUrl.href);

    try {
        const response = await fetch(cleanUrl.href);
        if (!response.ok) {
            console.error(`Failed to fetch ${cleanUrl.href}: ${response.statusText}`);
            return;
        }

        const html = await response.text();
        const { window } = new JSDOM(html);
        const document = window.document;

        const title = extractTitle(document, cleanUrl.href);
        const mainContent = extractMainContent(document);

        if (!mainContent) {
            return;
        }

        pagesToIngest.push({
            title,
            url: unwrapWaybackUrl(cleanUrl.href),
            language: 'en',
            content: mainContent
        });

        console.log(`[${pagesToIngest.length}/${MAX_PAGES}] Indexed: ${cleanUrl.href}`);

        if (depth >= MAX_DEPTH) {
            return;
        }

        const links = Array.from(document.querySelectorAll('a'))
            .map(link => link.getAttribute('href'))
            .filter(Boolean);

        for (let href of links) {
            if (pagesToIngest.length >= MAX_PAGES) {
                break;
            }

            if (href && href.startsWith('/wiki/') && !href.includes(':')) {
                const resolvedUrl = new URL(href, cleanUrl.origin).href;
                const normalizedUrl = resolvedUrl.split('#')[0];
                
                if (!visitedUrls.has(normalizedUrl)) {
                    await delay(CRAWL_DELAY_MS);
                    await crawlPage(normalizedUrl, depth + 1);
                }
            }
        }
    } catch (error) {
        console.error(`Error crawling ${cleanUrl.href}:`, error);
    }
}

async function sendBatchToBackend() {
    if (!ENABLE_INGEST) {
        console.log('[INFO] ENABLE_INGEST is false, skipping upload step');
        return;
    }

    if (pagesToIngest.length === 0) {
        console.log('[INFO] No pages collected, nothing to upload');
        return;
    }

    const url = new URL(INGEST_ENDPOINT, BACKEND_URL);
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Crawler-Key': CRAWLER_KEY
        },
        body: JSON.stringify(pagesToIngest)
    });

    if (!response.ok) {
        const errorBody = await response.text();
        throw new Error(`ingest endpoint failed: ${response.status} ${errorBody}`);
    }

    console.log(`[INFO] Uploaded ${pagesToIngest.length} pages to backend`);
}

async function main() {
    const seedURLs = await getSeedURLs();
    console.log(`[INFO] Using ${seedURLs.length} seed targets`);

    for (const seedURL of seedURLs) {
        if (pagesToIngest.length >= MAX_PAGES) {
            break;
        }

        await crawlPage(seedURL, 0);
    }

    console.log(`[INFO] Crawl complete. Collected ${pagesToIngest.length} pages`);
    await sendBatchToBackend();
}

main().catch((error) => {
    console.error('[FATAL] Crawler failed:', error.message);
    process.exit(1);
});