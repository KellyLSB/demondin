import React from 'react'
import { Form } from 'semantic-ui-react'

export default class ItemOption extends React.Component {
	render() {
		var option = this.props.option;

		if (option.valueType == "select") {
			var values;

			if (['key', 'text', 'value'] in option.values) {
				values = option.values;
			} else {
				values = option.values.map((v) => {
	 				return { key: v, text: v, value: v };
				});
			}			

			return (
				<Form.Field>
					<Form.Select label={option.key} fluid
											 placeholder={option.key} 
											 name={option.id} 
											 options={values} 
											 value={this.props.value}
											 onChange={this.props.onChange} />
				</Form.Field>
		)	};

		return (
			<Form.Field>
				<Form.Input label={option.key} fluid 
										placeholder={option.key} 
										name={option.id} 
										value={this.props.value}
										onChange={this.props.onChange} />
			</Form.Field>
	) };
}
