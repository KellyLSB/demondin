import React from 'react'
import { Form } from 'semantic-ui-react'

export default class ItemOption extends React.Component {
  constructor(props) {
    super(props);
  }

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
					<Form.Select name={option.key}
											 label={option.key} 
											 placeholder={option.key} 
											 options={values} />
				</Form.Field>
		)	};

		return (
			<Form.Field>
				<Form.Input name={option.key}
										label={option.key}
										placeholder={option.key} />
			</Form.Field>
	) };
}
