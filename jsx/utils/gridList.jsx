import React from 'react'
import { Grid } from 'semantic-ui-react'

export default class GridList extends React.Component {
	columns() {
		return this.props.columns || 2;
	}
	
	rows() {
		return Math.ceil(
			this.props.children.length / this.columns(),
		);
	}
	
	rand() { 
		return this.rand = this.rand || String.random(4); 
	}

	render() {
		var grid = [];

		for(var i = 0; this.rows() > i; i++) {
			var is = i * this.columns();
			
			grid.push(new Grid.Row( {
				key: `r${this.rand()}-${i}`,
				children: this.props.children.slice(
					is, is + this.columns(),
				).map((child, i) => new Grid.Column( {
					key: `c%${this.rand()}-${i}`,
					textAlign: i + 1 == this.columns() ?
						'right' : null,
					children: [ child ],
				} ) ),
			} ) );
		}
	
		return new Grid(
			Object.assign({}, this.props, { children: grid }),
		);
	}
}
