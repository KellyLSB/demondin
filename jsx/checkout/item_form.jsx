import React from 'react'
import { Mutation } from "react-apollo";
import gql from "graphql-tag";

import { Divider, Form, Button } from 'semantic-ui-react'

import ItemOption from './item_option'
import GridList from '../utils/gridList';
import FormHelper from '../utils/formHelper';

export default class ItemForm extends FormHelper {
	constructor(props) {
		super(props);

		this.onSubmit = this.onSubmit.bind(this);
	}

	onSubmit(e, updateInvoice) {
		e.preventDefault()
		
		var input = { variables: {
			input: { 
				items: [{
					itemID: this.props.item,
					itemPriceID: this.props.price,
					options: this.mapValues((name, value) => {
						return {
							itemOptionTypeID: name,
							values: value
						};
					})
				}] 
			}
		} };
		
		console.log(input);

		updateInvoice(input);
	}

	render() {
		return (
			 <Mutation mutation={gql`
				mutation activeInvoice($input: NewInvoice!) {
					activeInvoice(input: $input) {
						id
						items {
							id
							options {
								id
								itemOptionType {
									id
									key
								}
								values
							}
						}
					}
				}
			`}>
				{(updateInvoice) => (
					<Form onSubmit={(e) => this.onSubmit(e, updateInvoice)}>
						<GridList columns={2}>
							{this.props.options.map((option) => 
								<ItemOption key={option.id} option={option}
														getValue={this.getValue(option.key)}
														onChange={this.onChange} />
							) }
						</GridList>
						<Divider hidden />
						<Button type='submit'>Purchase</Button>
					</Form>
				) }
			</Mutation>
	) };
}
