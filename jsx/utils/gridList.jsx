import React from 'react'
import { Grid } from 'semantic-ui-react'

export default class GridList extends React.Component {
	columns() {
		return this.props.columns || 2;
	}

	render() {
		const rand = String.random(4);
		const rows = Math.ceil(
			this.props.children.length / this.columns(),
		);
		
		var grid = [];

		for(var i = 0; rows > i; i++) {
			var is = i * this.columns();
			
			grid.push(new Grid.Row( {
				key: `r${rand}-${i}`,
				children: this.props.children.slice(
					is, is + this.columns(),
				).map((child, i) => new Grid.Column( {
					key: `c%${rand}-${i}`,
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
