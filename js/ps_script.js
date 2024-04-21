function loadJsonParser(){
    var jsonParser = new File("C:\\Users\\teteu\\OneDrive\\Documentos\\Coding\\ps-automate\\js\\json2.js");
    jsonParser.open("r");
    var parser = jsonParser.read();
    eval(parser);
}

function loadParameters(){
    var jsonFile = new File("C:\\Users\\teteu\\OneDrive\\Documentos\\Coding\\ps-automate\\js\\parameters.json");
    jsonFile.open('r');
    var data = jsonFile.read();
    jsonFile.close();
    var castData = JSON.parse(data); 

    return castData
}

function getLayers(PSInstance, fileLayers) {
    var layers = {};
    
    for (var i = 0; i < fileLayers.length; i++) {
        var str = fileLayers[i].Name;
        layers[str] = null;
    }

    for (var j = 0; j < PSInstance.layers.length; j++) {
        var layer = PSInstance.layers[j];

        for (item in fileLayers) {
            if(fileLayers[item].Name === layer.name) {
                layers[fileLayers[item].Name] = layer
                continue
            }
        }
    }

    return layers

}

function getImage(caminhoImagem) {
    var novaImagem = new File(caminhoImagem);
    var novaImagemDocumento = app.open(novaImagem);
    novaImagemDocumento.selection.selectAll();
    novaImagemDocumento.selection.copy();
    novaImagemDocumento.close(SaveOptions.DONOTSAVECHANGES);
}

function pasteImage(targetLayer) {
    app.preferences.rulerUnits = Units.PIXELS;
    app.activeDocument = targetLayer.parent;
    app.activeDocument.activeLayer = targetLayer;
    app.activeDocument.paste();
}

function resizeImg(layer, maxWidth, maxHeight) {
    var width = layer.bounds[2] - layer.bounds[0];
    var height = layer.bounds[3] - layer.bounds[1];
    var scale = Math.min(maxWidth / width, maxHeight / height);
    layer.resize(scale * 100, scale * 100, AnchorPosition.MIDDLECENTER);
}

function changeText(layer, newText) {
    layer.textItem.contents = newText
}

function changeValue(layer, newContent, layerType) {
    if(layerType === "Text") {
        changeText(layer, newContent)
        return
    }
    changeImageBuilder(layer, newContent)
}

function positionImage(targetLayer) {
    var pastedLayer = app.activeDocument.activeLayer;
    var pastedCenterX = (pastedLayer.bounds[0] + pastedLayer.bounds[2]) / 2;
    var pastedCenterY = (pastedLayer.bounds[1] + pastedLayer.bounds[3]) / 2;

    var targetCenterX = (targetLayer.bounds[0] + targetLayer.bounds[2]) / 2;
    var targetCenterY = (targetLayer.bounds[1] + targetLayer.bounds[3]) / 2;

    var maxWidth = targetLayer.bounds[2] - targetLayer.bounds[0];
    var maxHeight = targetLayer.bounds[3] - targetLayer.bounds[1];

    var xOffset = targetCenterX - pastedCenterX;
    var yOffset = targetCenterY - pastedCenterY;

    pastedLayer.translate(xOffset, yOffset);
    targetLayer.visible = false;
    resizeImg(pastedLayer, maxWidth, maxHeight)
}

function changeImageBuilder(targetLayer, imgPath) {
    getImage(imgPath)
    pasteImage(targetLayer)
    positionImage(targetLayer)
}

function exportAsPNG(doc, finalName) {
    var exportOptions = new ExportOptionsSaveForWeb();
    exportOptions.format = SaveDocumentType.PNG;
    exportOptions.PNG8 = false; // Set to true for PNG-8 format
    exportOptions.quality = 100; // Set PNG quality (0-100)
    doc.exportDocument(new File(finalName), ExportType.SAVEFORWEB, exportOptions);
}

function exportAsJPG(doc, finalName) {
    var exportOptions = new ExportOptionsSaveForWeb();
    exportOptions.format = SaveDocumentType.JPEG;
    exportOptions.quality = 100; // Set PNG quality (0-100)
    doc.exportDocument(new File(finalName), ExportType.SAVEFORWEB, exportOptions);
}

function exportAsPSD(doc, finalName) {
    var saveOptions = new PhotoshopSaveOptions();
    doc.saveAs(new File(finalName), saveOptions, true);
}

var ExportTypes = {
    0: "psd",
    1: "png",
    2: "jpg"
}

var exportFunctions = {
    "psd": exportAsPSD,
    "png": exportAsPNG,
    "jpg": exportAsJPG
};

function exportFile(doc, finalPath, type) {
    var exportType = ExportTypes[type];
    var finalPathWithName = finalPath + '.' + ExportTypes[type]

    exportFunctions[exportType](doc, finalPathWithName);    
}

loadJsonParser()
var parameters = loadParameters()
var items = parameters.Data
var exportDir = parameters.ExportDir
var prefix = parameters.PrefixNameForFile
var exportTypes = parameters.ExportTypes
var psdTemplate = parameters.PSDTemplate

for (index in items) {
    var exportPath = exportDir + "\\" + prefix + '-' + index
    app.open(new File(psdTemplate));
    var documento = app.activeDocument;
    var layers = getLayers(documento, parameters.Layers)

    for (i in layers) {
        for(j in parameters.Layers) {
            if(parameters.Layers[j].Name === layers[i].name) {
                var layerType = parameters.Layers[j].Type
            }
        }
        var newContent = items[index][layers[i].name]
        changeValue(layers[i], newContent, layerType)
        
    }

    for (k in exportTypes) {
        exportFile(documento, exportPath, exportTypes[k])
    }

    documento.close(SaveOptions.DONOTSAVECHANGES);
    
}
