const MEGABYTE = 1024 * 1024;

function memIntensive(level, memSize) {
    const available_memory = parseInt(memSize, 10);
    const amountInMB = available_memory - (available_memory / 10) * (4 - level);
    console.log(amountInMB);
    Buffer.alloc(amountInMB * MEGABYTE, 'a');
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

function runTest(intensityLevel, memSize){
    memIntensive(intensityLevel, memSize);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.memstress = (args) => {
    const startTime = process.hrtime();
    const params = getParameters(args);
    if (params.error) {
        return {
            statusCode: 200,
            headers: { 'Content-Type': 'application/json' },
            body: { error: params.error },
        };
    }

    runTest(params, args.function_memory_size);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        statusCode: 200,
        headers: {
            'Content-Type': 'application/json',
        },
        body: { reused, duration },
    }
}
