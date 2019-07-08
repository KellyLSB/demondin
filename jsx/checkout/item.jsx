import React from 'react';

import { Subscription } from "react-apollo";
import gql from "graphql-tag";

import { Divider, Form,
				  Header, Segment,
					Button, Icon, Label
} from 'semantic-ui-react';

import ItemForm from './item_form';
import GridList from '../utils/gridList';

export default class Item extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			hideForm: true,
		};

		this.closeForm = this.closeForm.bind(this);
		this.onToggleForm = this.onToggleForm.bind(this);
	}
	
	closeForm() {
		this.setState((state) => {
			state.hideForm = true;
			return state;
		});
	}

	onToggleForm() {
		this.setState((state) => {
			state.hideForm = !state.hideForm;
			return state;
		});
	}

	toDollars(cents) {
		return cents.toDollars(false);
	}

	toCents(dollars) {
		return dollars.toCents();
	}

	currentPrice() {
		return this.props.item.prices[0];
	}

	printPrice() {
		var price = this.currentPrice();

		if (price) {		
			if('taxable' in price && price.taxable) {
				return this.toDollars(price.price) + " + Tax"
			}

			return this.toDollars(price.price)
		}

		return "Undefined"
	}

	render() {
		return <React.Fragment>
			<Header attached='top'>
				{this.props.item.name}
			</Header>

			<Segment attached>
				{this.props.item.description}
				
				<Divider hidden />

				<GridList columns={2}>
					<Label ribbon>
						<Icon name='dollar sign' />
						{ this.printPrice() }
					</Label>

					<InvoiceCharged not>
						<Button onClick={this.onToggleForm} 
										primary icon labelPosition='left'>
							<Icon name='shop' />
							Customize
						</Button>
					</InvoiceCharged>
				</GridList>
			</Segment>

			<InvoiceCharged not>
				<ItemForm item={this.props.item.id} hideForm={this.state.hideForm}
									price={this.currentPrice().id} 
									options={this.props.item.options} />
			</InvoiceCharged>
		</React.Fragment>;
	}
}

class InvoiceCharged extends React.Component {
	render() {
		return <Subscription subscription={gql`
			subscription InvoiceUpdated {
				invoiceUpdated {
					id
					stripeChargeID
				}
			}`}>
			{({ data, loading }) => {
				var invoice = data ? data.invoiceUpdated : false;
				if(loading || !invoice || invoice.stripeChargeID) {
					if(!this.props.not) return this.props.children;
					return null;
				}
				
				if(this.props.not) return this.props.children;
				return null;
			} }
		</Subscription>;
	}
}
