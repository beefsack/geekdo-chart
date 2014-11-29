package main

import (
	"html/template"
)

var rawTemplate = `<html>
<head>
<script src="http://www.amcharts.com/lib/3/amcharts.js"></script>
<script src="http://www.amcharts.com/lib/3/serial.js"></script>
<script>
document.addEventListener("DOMContentLoaded", function(event) {
	AmCharts.makeChart("chartdiv", {
		type: "serial",
		pathToImages: "http://www.amcharts.com/lib/3/images/",
		dataProvider: {{.DataProvider}},
		valueAxes: [{
			logarithmic: true,
			reversed: true,
			minimum: 1
		}],
		legend: {},
		graphs: {{.Graphs}},
		chartScrollbar: {},
		chartCursor: {
			cursorPosition: "mouse"
		},
		dataDateFormat: "YYYY-MM-DD",
		categoryField: "date",
		categoryAxis: {
			parseDates: true
		}
	});
});
</script>
<style>
#chartdiv {
	width: 100%;
	height: 700px;
}
</style>
</head>
<body>
<div id="chartdiv"></div>
</body>
</html>`
var parsedTemplate = template.Must(template.New("template").Parse(rawTemplate))
