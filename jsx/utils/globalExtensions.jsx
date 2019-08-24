// Extend some global class prototypes
String.random = function(len = 4) {
	return Math.random().toString(36).substr(2, len);
}

String.prototype.toCents = function() {
	return Math.round(100 * parseFloat(
		this.replace(/[$,]/g, '')
	));
}

Number.prototype.toDollars = function(sign = true) {
	return Intl.NumberFormat('en-US', {
		style: 'currency', 
		currency: 'USD',
	}).formatToParts(this / 100).map(({ type, value }) => {
		switch(type) {
		case 'currency': return sign ? value : '';
		default: return value;
		}
	}).join('');
}
