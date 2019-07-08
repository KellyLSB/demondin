import React from 'react';

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

		this.onToggleForm = this.onToggleForm.bind(this)
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
					<Button primary icon labelPosition='left'
						onClick={this.onToggleForm}>
						<Icon name='shop' />
						Customize
					</Button>
				</GridList>
			</Segment>

			{ this.state.hideForm ? null : (
				<Segment secondary attached='bottom'>
					<ItemForm item={this.props.item.id} 
										price={this.currentPrice().id} 
										options={this.props.item.options} />
				</Segment>
			) }
		</React.Fragment>;
	}
}
