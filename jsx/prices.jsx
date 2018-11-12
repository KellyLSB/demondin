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

        this.addPrice = this.addPrice.bind(this)
    }

    addPrice() {
        console.log("POKEY")
        this.setState((state) => {
            state.data.push({
              ItemID: this.props.item
            })
            
            return state
        })
    }

    render() {
        return (
            <React.Fragment>
                {this.state.data.map((price, altID) =>
                    <Price key={price.ID ? price.ID : altID}
                      item={this.props.item}
                      data={price} />
                )}
                <input type="button" value="Add Another Price"
                    onClick={this.addPrice} />
            </React.Fragment>
        )
    }
}
