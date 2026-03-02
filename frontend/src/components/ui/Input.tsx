/**
 * Input — adapter wrapper around @iit/ui Input.
 */
import React from 'react'
import { Input as IITInput } from '@iit/ui'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  helperText?: string
  icon?: React.ReactNode
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ icon, ...props }, ref) => {
    return (
      <IITInput
        ref={ref}
        leftElement={icon}
        {...props}
      />
    )
  }
)

Input.displayName = 'Input'

export default Input
