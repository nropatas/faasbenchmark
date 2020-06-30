const proc = require('child_process');
const url = require('url');

const PATH = '/tmp/faastest';

function ioIntensive(baseNumber) {
    const amountInMB = 10 ** (baseNumber - 1);
    const out = proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    if (out.status !== 0)
        return out.stderr.toString();
}

function getQueryParams(context) {
    const urlParts = url.parse(context.request.url, true);
    return urlParts.query;
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel) {
    return ioIntensive(intensityLevel)
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

module.exports = async (context) => {
    const query = getQueryParams(context);

    const startTime = process.hrtime();
    const params = getParameters(query);
    if (params.error) {
        return {
            status: 200,
            body: { error: params.error },
        };
    }

    const error = runTest(params);
    if (error) {
        return {
            status: 200,
            body: { error: error },
        };
    }

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        status: 200,
        body: { reused, duration },
    };
};
