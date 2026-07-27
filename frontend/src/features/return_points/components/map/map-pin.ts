export const MAP_PIN_IMAGE_PIXEL_RATIO = 2;

export const getMapPinImageId = (
    bgColor: string,
    textColor: string,
    shortName: string,
): string => {
    const imageVariant = [bgColor, textColor, shortName]
        .map((value) => value.trim().toLowerCase())
        .join("\u0000");
    let hash = 0;

    for (const character of imageVariant) {
        hash = (hash * 31 + character.charCodeAt(0)) >>> 0;
    }

    return `map-pin-${hash.toString(36)}`;
};

export const createMapPinImageData = (
    bgColor: string,
    textColor: string,
    shortName: string,
): ImageData => {
    const canvas = document.createElement("canvas");
    const imageSize = 64;
    canvas.width = imageSize * MAP_PIN_IMAGE_PIXEL_RATIO;
    canvas.height = imageSize * MAP_PIN_IMAGE_PIXEL_RATIO;

    const context = canvas.getContext("2d");
    if (!context) throw new Error("Nie udało się utworzyć grafiki markera mapy.");

    context.scale(MAP_PIN_IMAGE_PIXEL_RATIO, MAP_PIN_IMAGE_PIXEL_RATIO);
    context.translate(imageSize / 2, imageSize / 2);
    context.rotate(Math.PI / 4);
    context.beginPath();
    context.roundRect(-20, -20, 40, 40, [20, 20, 20, 6]);
    context.fillStyle = bgColor.trim() || "#2f8f46";
    context.strokeStyle = "#ffffff";
    context.lineWidth = 3;
    context.shadowColor = "rgba(14, 67, 40, .28)";
    context.shadowBlur = 6;
    context.shadowOffsetY = 2;
    context.fill();
    context.shadowColor = "transparent";
    context.stroke();

    context.rotate(-Math.PI / 4);
    context.fillStyle = textColor.trim() || "#ffffff";
    context.font = '800 12px "Bricolage Grotesque", ui-sans-serif, system-ui, sans-serif';
    context.textAlign = "center";
    context.textBaseline = "middle";
    context.fillText(shortName.trim(), 0, 0, 30);

    return context.getImageData(0, 0, canvas.width, canvas.height);
};