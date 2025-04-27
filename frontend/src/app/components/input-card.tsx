import { useState } from "react";
import { CalculationResult } from "../model";
import { Result } from "./result";
import styles from "../page.module.css"





export function InputCard({onCalculate, currentResult}: 
    {onCalculate: (input: string) => Promise<void>, currentResult?: CalculationResult | null}) {
  
    const [calculationInputValue, setCalculationInputValue] = useState<string>("");
  
    function onKeyUp(e: any) {
      if (e && e.key == 'Enter') {
        onCalculate(calculationInputValue)
      }
    }
  
    return (
      <>
       <div className={`${styles.inputSection} ${styles.gridCard}`}>
          <h2>Input</h2>
          <div className={styles.inputMain}>
            <div className={styles.inputForm}>
              <input placeholder="1+(3+4)" 
                value={calculationInputValue} 
                onChange={ e => setCalculationInputValue(e.currentTarget.value)} 
                onKeyUp={(e) => onKeyUp(e)}></input>
              <button onClick={() => onCalculate(calculationInputValue)}>Calculate!</button>
  
            </div>  
            { currentResult ? <Result calculationResult={currentResult}></Result> : null}
          </div>
      </div>
      </>
    )
}
  