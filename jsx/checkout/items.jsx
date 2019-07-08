import React from 'react'
import { Query } from "react-apollo";
import gql from "graphql-tag";

import Item from './item'

import { Grid } from 'semantic-ui-react'

export default class Items extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			data: []
		};
	}

	render() {
		return (
			<Grid.Row>
				<Query query={gql`{
					items {
						id
						name
						description
						prices {
							id
							price
							beforeDate
							afterDate
						}
						options {
							id
							key
							valueType
							values
						}
					}
				}`}>
					{({ loading, error, data }) => {
						if (loading) return <p>Loading...</p>;
						if (error) return <p>Error :(</p>;

						return data.items.map((item) => 
							<Item key={item.id} item={item} />
						);
					}}
				</Query>
			</Grid.Row>
		)
	}
}
