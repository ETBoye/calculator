import styles from "../page.module.css"

export function InfoCard() {
    return (
      <div className={styles.gridCard}>
        <h2> Info </h2>
        <p>This is a calculator for computing mathematical expressions with integers.</p>
        You can input
        <ul>
          <li> Integer constants, eg. 3, 5, 0, -0, -432</li>
          <li> The usual 4 <strong>binary</strong> operations, eg. 3+4, 4-2, 3*2, 5/2</li>
          <li> Any expressions of the above using brackets, eg. (4+3)*2/(5-2)</li>
        </ul>
        You <strong>can not</strong> input 
        <ul>
          <li> The unary operator '-' when not used in the context of a negative constant, eg. -(2+3), --5</li>
          <li> The unary operator '+' in any case, eg. +3</li>
          <li> Any decimal representation of numbers, eg. 3.14, 100,00</li>
        </ul>
  
        <p>The calculator respects the order of operations, and multiplication and division is evaluated left to right. For instance,</p>
        <p>3*4/2*5 = ((3*4)/2)*5</p>
  
        <p>The calculator computes the result accurately as a rational number - however, the input and result is stored in a database which restricts the length of the input and output</p>
      </div>
    )
}