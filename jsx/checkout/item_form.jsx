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

	onChange(e, key, data) {
		this.setState((state) => {
			state.values[key] = data.value;
			return state
		} );

		console.log(e, key, data);
		debugger
	}

	onSubmit(e) {
		e.preventDefault()
	}

  render() {
		var option = this.props.option;

		return (
			<Form onSubmit={this.onSubmit}>
			  {this.state.data.options.map((option) => 
				  <ItemOption key={option.key} option={option}
											value={this.state[option.key]}
                      onChange={(e,d) => this.onChange(e,option.key,d)} />
				) }
				<Button type='submit'>Purchase</Button>
			</Form>
	) };
}
