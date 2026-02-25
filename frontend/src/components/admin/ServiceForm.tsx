import { useEffect } from 'react'
import { useForm, useFieldArray } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Plus, Trash2 } from 'lucide-react'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import type { Service, ServiceCreateRequest, ServiceUpdateRequest } from '@/types/auth'

const schema = z.object({
  name: z.string().min(2, 'Minimo 2 caracteres'),
  slug: z.string().min(2, 'Minimo 2 caracteres').regex(/^[a-z0-9-]+$/, 'Apenas letras minusculas, numeros e hifens'),
  description: z.string().optional(),
  redirect_urls: z.array(z.object({ url: z.string().url('URL invalida') })).optional(),
})

type FormValues = z.infer<typeof schema>

interface ServiceFormProps {
  service?: Service | null
  onSubmit: (data: ServiceCreateRequest | ServiceUpdateRequest) => Promise<void>
  onCancel: () => void
}

function slugify(value: string): string {
  return value
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9\s-]/g, '')
    .trim()
    .replace(/\s+/g, '-')
}

export function ServiceForm({ service, onSubmit, onCancel }: ServiceFormProps) {
  const isEditing = !!service

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    control,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: isEditing
      ? {
          name: service.name,
          slug: service.slug,
          description: service.description || '',
          redirect_urls: (service.redirect_urls || []).map((url) => ({ url })),
        }
      : {
          name: '',
          slug: '',
          description: '',
          redirect_urls: [],
        },
  })

  const { fields, append, remove } = useFieldArray({ control, name: 'redirect_urls' })

  const nameValue = watch('name')

  useEffect(() => {
    if (!isEditing) {
      setValue('slug', slugify(nameValue || ''))
    }
  }, [nameValue, isEditing, setValue])

  const handleFormSubmit = async (values: FormValues) => {
    const payload: ServiceCreateRequest | ServiceUpdateRequest = {
      name: values.name,
      slug: values.slug,
      description: values.description || undefined,
      redirect_urls: values.redirect_urls?.map((r) => r.url).filter(Boolean) || [],
    }
    await onSubmit(payload)
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          {...register('name')}
          label="Nome"
          placeholder="Meu Servico"
          error={errors.name?.message}
          required
        />
        <Input
          {...register('slug')}
          label="Slug"
          placeholder="meu-servico"
          helperText="Identificador unico, gerado automaticamente"
          error={errors.slug?.message}
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">
          Descricao
        </label>
        <textarea
          {...register('description')}
          placeholder="Descricao do servico (opcional)"
          rows={2}
          className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors resize-none"
        />
      </div>

      <div>
        <div className="flex items-center justify-between mb-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-200">
            URLs de redirecionamento
          </label>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            icon={<Plus className="w-3.5 h-3.5" />}
            onClick={() => append({ url: '' })}
          >
            Adicionar URL
          </Button>
        </div>

        <div className="space-y-2">
          {fields.map((field, index) => (
            <div key={field.id} className="flex items-start gap-2">
              <div className="flex-1">
                <Input
                  {...register(`redirect_urls.${index}.url`)}
                  placeholder="https://meuapp.com/callback"
                  error={errors.redirect_urls?.[index]?.url?.message}
                />
              </div>
              <button
                type="button"
                onClick={() => remove(index)}
                className="mt-2 p-1.5 text-gray-400 hover:text-danger-600 dark:hover:text-danger-400 transition-colors"
                title="Remover"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </div>
          ))}
          {fields.length === 0 && (
            <p className="text-sm text-gray-400 dark:text-gray-500 italic">
              Nenhuma URL adicionada
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3 pt-2">
        <Button type="submit" loading={isSubmitting}>
          {isEditing ? 'Salvar alteracoes' : 'Criar servico'}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}

export default ServiceForm
