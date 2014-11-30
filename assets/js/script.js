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

function initChartPage(graphs, dataProvider) {
	makeChart("chartdiv", graphs, dataProvider);
	$('#thingChooser').select2({
		width: '50%',
		minimumInputLength: 2,
		multiple: true,
		ajax: {
			url: '/search',
			dataType: 'json',
			quietMillis: 250,
			data: function(term, page) {
				return {
					query: term
				};
			},
			results: function(data, page) {
				return {
					results: $.map(data, function(d) {
						return {
							id: d.type + ':' + d.id,
							text: d.name
						};
					})
				};
			},
			cache: true
		},
		initSelection: function(element, callback) {
			$(element).val('');
			callback($.map(graphs, function(g) {
				return {
					id: g.valueField,
					text: g.title
				};
			}));
		}
	});
}

function submitChooseThingsForm(el) {
	var val = $('#thingChooser').val();
	if (!val) {
		alert('Please select at least one item');
		return;
	}
	document.location = '/' + val;
}
