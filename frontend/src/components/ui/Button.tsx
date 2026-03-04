/**
 * Button — adapter wrapper around @iit/ui Button.
 * Preserves the existing API used across the app while delegating to the design system.
 */
import React from 'react'
import { Button as IITButton, type ButtonProps as IITButtonProps } from '@iit/ui'
import { cn } from '@/lib/utils'

interface ButtonProps extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, 'onDrag' | 'onDragStart' | 'onDragEnd' | 'onAnimationStart'> {
  variant?: 'primary' | 'secondary' | 'success' | 'danger' | 'ghost' | 'outline' | 'accent' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  isLoading?: boolean
  icon?: React.ReactNode
  leftIcon?: React.ReactNode
  rightIcon?: React.ReactNode
  fullWidth?: boolean
}

/** Map local variant names to @iit/ui variants */
function mapVariant(v?: ButtonProps['variant']): IITButtonProps['variant'] {
  if (v === 'danger') return 'destructive'
  if (v === 'success') return 'accent'   // closest in IIT palette
  if (v === 'outline') return 'secondary'
  return (v as IITButtonProps['variant']) ?? 'primary'
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      children,
      variant = 'primary',
      size = 'md',
      loading,
      isLoading,
      icon,
      leftIcon,
      rightIcon,
      fullWidth = false,
      className,
      ...props
    },
    ref
  ) => {
    const resolvedLoading = loading || isLoading
    const resolvedLeftIcon = leftIcon || icon

    return (
      <IITButton
        ref={ref}
        variant={mapVariant(variant)}
        size={size}
        isLoading={resolvedLoading}
        leftIcon={resolvedLeftIcon}
        rightIcon={rightIcon}
        className={cn(fullWidth && 'w-full', className)}
        {...props}
      >
        {children}
      </IITButton>
    )
  }
)

Button.displayName = 'Button'

export { Button }
export default Button
