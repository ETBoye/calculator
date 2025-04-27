import { MathJax } from "better-react-mathjax";
import { CalculationResult } from "../model";
import styles from "../page.module.css"

export function Result({calculationResult} : {calculationResult: CalculationResult}) {

    function  getMathJaxString() {
  
      if (!calculationResult.result) {
        return "";
      }
  
      if (calculationResult.result.denom === "1") {
        return `\\[${calculationResult.result.num}\\]`
      }
      
      return  `\\[\\frac{${calculationResult.result.num}}{${calculationResult.result.denom}}  \\approx ${calculationResult.result.estimate}\\]`;
    }
    return (
        <div className={styles.center}>
          {
            calculationResult.errorId ? <p>{calculationResult.errorId}</p> 
            : <MathJax inline dynamic>{getMathJaxString()}</MathJax> 
          }
          </div>
    )
  }
  