import React from 'react'
import Price from './price'

export default class Prices extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            data: []
        };

        if ('data' in props) {
            this.state.data = props.data;
        }

        this.addPrice =  this.addPrice.bind(this)
    }

    addPrice() {
        this.setState((state) => {
            state.data += {}
        })
    }

    render() {
        return (
            <React.Fragment>
                {this.state.data.map((price) =>
                    <Price key={price.ID} data={price} />
                )}
                <input type="button" value="Add Another Price"
                    onClick={this.addPrice} />
            </React.Fragment>
        )
    }
}
