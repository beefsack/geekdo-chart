function makeChart(el, graphs, chartData) {
	AmCharts.makeChart(el, {
		type: "serial",
		pathToImages: "http://www.amcharts.com/lib/3/images/",
		dataProvider: chartData,
		valueAxes: [{
			logarithmic: true,
			reversed: true,
			minimum: 1
		}],
		legend: {},
		graphs: graphs,
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
}
