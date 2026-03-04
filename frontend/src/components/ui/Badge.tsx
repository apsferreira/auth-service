/**
 * Badge — adapter wrapper around @iit/ui Badge.
 * Maps local variant names to @iit/ui variant names.
 */
import React from 'react'
import { Badge as IITBadge, type BadgeProps as IITBadgeProps } from '@iit/ui'

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info' | 'outline' | 'neutral' | 'error' | 'brand'
  size?: 'sm' | 'md' | 'lg'
  dot?: boolean
}

function mapVariant(v?: BadgeProps['variant']): IITBadgeProps['variant'] {
  if (v === 'default' || v === 'outline') return 'neutral'
  if (v === 'danger') return 'error'
  return (v as IITBadgeProps['variant']) ?? 'neutral'
}

const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({ variant, dot, children, ...props }, ref) => {
    return (
      <IITBadge ref={ref} variant={mapVariant(variant)} dot={dot} {...props}>
        {children}
      </IITBadge>
    )
  }
)

Badge.displayName = 'Badge'

export { Badge }
export default Badge
