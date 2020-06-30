const proc = require('child_process');

const PATH = '/tmp/faastest';

function ioIntensive(baseNumber) {
    const amountInMB = 10 ** (baseNumber - 1);
    const out = proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    if (out.status !== 0)
        return out.stderr.toString();
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.query.level ? parseInt(event.query.level) : null;
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

module.exports = async (event, context) => {
    const startTime = process.hrtime();
    const params = getParameters(event);
    if (params.error) {
        return context.status(200).succeed({
            error: params.error,
        });
    }

    const error = runTest(params);
    if (error) {
        return context.status(200).succeed({
            error: error,
        });
    }

    const reused = isWarm();
    const duration = getDuration(startTime);

    return context.status(200).succeed({
        reused: reused,
        duration: duration,
    });
}
