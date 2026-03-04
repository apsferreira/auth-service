/**
 * Select — adapter wrapper around @iit/ui Select.
 */
import React from 'react'
import { Select as IITSelect } from '@iit/ui'

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string
  error?: string
  helperText?: string
  options: { value: string; label: string }[]
  placeholder?: string
}

const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  (props, ref) => {
    return <IITSelect ref={ref} {...props} />
  }
)

Select.displayName = 'Select'

export default Select
