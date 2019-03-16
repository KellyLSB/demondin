import React from 'react'
import { Mutation } from "react-apollo";
import gql from "graphql-tag";

import { Form, Button } from 'semantic-ui-react'

import ItemOption from './item_option'

export default class ItemForm extends React.Component {
  constructor(props) {
    super(props);

		this.state = {
			data: {},
			values: {},
		};

		if ('options' in props) {
      this.state.data.options = props.options;
    }

		this.onChange = this.onChange.bind(this);
		this.onSubmit = this.onSubmit.bind(this);
  }

	onChange(e, option, data) {
		this.setState((state) => {
			state.values[option.id] = data.value;
			return state
		} );

		console.log(e, key, data);
		debugger
	}

	onSubmit(e, updateInvoice) {
		e.preventDefault()

		var options = this.state.values.map((key, value) => {
			optionTypeID: value,
			values: value
		})

		updateInvoice({ variables: {
			input: { options: options }
		} })
	}

  render() {
		var option = this.props.option;

		return (
			 <Mutation mutation={gql`
				mutation updateInvoice($id: ID!, $input: NewInvoice!) {
					updateInvoice(id: $id, input: $input) {
						id
						items {
							id
							options {
								id
								optionType {
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
						{this.state.data.options.map((option) => 
							<ItemOption key={option.key} option={option}
													value={this.state[option.key]}
				                  onChange={(e, d) => this.onChange(e, option, d)} />
						) }
						<Button type='submit'>Purchase</Button>
					</Form>
				) }				
			</Mutation>
	) };
}
