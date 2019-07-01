import React from 'react';

import { Mutation } from "react-apollo";
import gql from "graphql-tag";

import {
	Divider, Segment, Header, Label, Form, Message, Button,
} from 'semantic-ui-react';

import GridList from '../utils/gridList';
import StateErrors from '../utils/stateErrors';

const MUTATE_ACTIVE_INVOICE_ID = gql(`
mutation activeInvoice($input: NewInvoice!) {
	activeInvoice(input: $input) {
		id
	}
}`);

export default class CartItem extends StateErrors {
	constructor(props) {
		super(props)
		
		this.onRemove = this.onRemove.bind(this);
	}
		
	onRemove(e, item, updateInvoice) {
		e.preventDefault();
		
		updateInvoice({ variables: {
			input: {
				items: [{
					id: item.id,
					itemID: item.itemID,
					itemPriceID: item.itemPriceID,
					options: [],
					remove: true,
				}],
			},
		} }).catch((error) => {
			this.onError("removeCartItem", error);
		});
	}

	render() {
		const item = this.props.item;
		return <Segment key={item.id} attached>
			<Header as='h4' dividing>
				{ item.item ? 
					item.item.name : "<Item>"
				}
			</Header>
			<Label attached='top right' size='mini'>
				<Mutation mutation={MUTATE_ACTIVE_INVOICE_ID}>
					{ (updateInvoice) => (
						<Form error={ this.hasError() }
							onSubmit={ (e) => this.onRemove(e, item, updateInvoice) }>
							<Message error header='Error removing from cart'
								list={this.getErrors()}
							/>
							<Button icon='close' type='submit' size='mini' circular />
						</Form>
					) }
				</Mutation>
			</Label>

			<GridList columns={2}>
				{ item.options.map((option, i) =>
					<React.Fragment key={option.itemOptionType ? 
						option.itemOptionType.id : `option-${i}`
					}>
						<Header>
							{ option.itemOptionType ? 
								option.itemOptionType.key : "<Option>"
							}
						</Header>
						{ option.values }
					</React.Fragment>
				) }
			</GridList>

			<Divider hidden />

			<Label ribbon='right'>
				{ item.itemPrice ? 
					item.itemPrice.price.toDollars() : "-.-"
				}
			</Label>
		</Segment>;
	}
}
