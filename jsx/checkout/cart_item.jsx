import React from 'react';

import { Mutation } from "react-apollo";
import gql from "graphql-tag";

import {
	Divider, Segment, Header, Label, Form, Message, Button,
} from 'semantic-ui-react';

import GridList from '../utils/gridList'

export default class CartItem extends React.Component {
	constructor(props) {
		super(props)
		
		this.state = {
			errors: [],
		};
		
		this.onError = this.onError.bind(this);
		this.hasError = this.hasError.bind(this);
		this.onRemove = this.onRemove.bind(this);
	}
	
	onError(name, msg = null) {
		this.setState((state) => {
			if (msg === false) {
				delete state.errors[name];
			} else {
				state.errors[name] = msg;
			}
			
			return state;
		} );
	}
	
	hasError(name = null) {
		if (name === null) {
			return Object.keys(this.state.errors).length > 0;
		}
		
		return this.state.errors.hasOwnProperty(name);
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
				<Mutation mutation={gql`
					mutation activeInvoice($input: NewInvoice!) {
						activeInvoice(input: $input) {
							id
						}
					}
				`}>
					{ (updateInvoice) => (
						<Form error={ this.hasError() }
							onSubmit={ (e) => this.onRemove(e, item, updateInvoice) }>
							<Message error header='Error removing from cart'
								list={Object.values(this.state.errors)}
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
