'use client'

import { motion, useScroll, useTransform, useInView } from 'framer-motion'
import { useRef, useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Sparkles, Zap, Download, Clock, Play, Check, ArrowRight, Video, Scissors, Brain } from 'lucide-react'

const AnimatedCounter = ({ target, duration = 2 }) => {
  const [count, setCount] = useState(0)
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true })

  useEffect(() => {
    if (!isInView) return
    
    let start = 0
    const end = target
    const incrementTime = (duration * 1000) / end
    const timer = setInterval(() => {
      start += 1
      setCount(start)
      if (start === end) clearInterval(timer)
    }, incrementTime)

    return () => clearInterval(timer)
  }, [isInView, target, duration])

  return <span ref={ref}>{count}</span>
}

const GalaxyBackground = () => {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none">
      {/* Stars */}
      {[...Array(200)].map((_, i) => (
        <motion.div
          key={`star-${i}`}
          className="absolute rounded-full"
          style={{
            width: Math.random() * 3 + 1,
            height: Math.random() * 3 + 1,
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            backgroundColor: ['#ffffff', '#a78bfa', '#60a5fa', '#fbbf24'][Math.floor(Math.random() * 4)]
          }}
          animate={{
            opacity: [0.2, 1, 0.2],
            scale: [1, 1.5, 1]
          }}
          transition={{
            duration: Math.random() * 3 + 2,
            repeat: Infinity,
            ease: 'easeInOut'
          }}
        />
      ))}
      
      {/* Shooting Stars */}
      {[...Array(3)].map((_, i) => (
        <motion.div
          key={`shooting-${i}`}
          className="absolute w-1 h-20 bg-gradient-to-b from-white to-transparent"
          style={{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 50}%`,
            transform: 'rotate(45deg)'
          }}
          animate={{
            x: [0, 1000],
            y: [0, 500],
            opacity: [0, 1, 0]
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
            repeatDelay: Math.random() * 10 + 5,
            ease: 'easeOut'
          }}
        />
      ))}
    </div>
  )
}

const HeroSection = () => {
  const ref = useRef(null)
  const { scrollYProgress } = useScroll({
    target: ref,
    offset: ['start start', 'end start']
  })
  
  const y = useTransform(scrollYProgress, [0, 1], ['0%', '50%'])
  const opacity = useTransform(scrollYProgress, [0, 1], [1, 0])

  return (
    <section ref={ref} className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Galaxy Background */}
      <motion.div 
        style={{ y }}
        className="absolute inset-0 z-0"
      >
        <div className="absolute inset-0 bg-gradient-to-b from-purple-900/40 via-indigo-900/30 to-black z-10" />
        <img 
          src="https://images.unsplash.com/photo-1557264337-e8a93017fe92?crop=entropy&cs=srgb&fm=jpg&ixid=M3w3NTY2NzR8MHwxfHNlYXJjaHwxfHxhYnN0cmFjdCUyMHRlY2h8ZW58MHx8fGJsdWV8MTc3MDg2Mjg2MXww&ixlib=rb-4.1.0&q=85"
          alt="Galaxy Background"
          className="w-full h-full object-cover opacity-70"
        />
      </motion.div>

      <GalaxyBackground />

      {/* Content */}
      <motion.div 
        style={{ opacity }}
        className="relative z-20 container mx-auto px-6 text-center"
      >
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/10 border border-purple-500/30 backdrop-blur-md mb-8"
        >
          <Sparkles className="w-4 h-4 text-purple-300" />
          <span className="text-sm text-purple-200">AI-Powered Clipping Video</span>
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          className="text-6xl md:text-8xl font-bold mb-6 bg-gradient-to-b from-white via-purple-200 to-purple-400 bg-clip-text text-transparent leading-tight"
        >
          Ubah Video Menjadi<br />Klip Viral
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          className="text-xl md:text-2xl text-gray-300 mb-12 max-w-3xl mx-auto leading-relaxed"
        >
          Biarkan AI secara otomatis mendeteksi dan memotong momen terbaik dari video YouTube Anda.
          <span className="text-white"> Tanpa perlu keahlian editing.</span>
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.8 }}
          className="flex flex-col sm:flex-row gap-4 justify-center items-center"
        >
          <Button 
            size="lg" 
            className="group relative px-8 py-6 text-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 border-0 shadow-lg shadow-purple-500/50 hover:shadow-purple-500/70 transition-all duration-300"
          >
            <span className="relative z-10 flex items-center gap-2">
              Mulai Gratis
              <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </span>
            <div className="absolute inset-0 bg-gradient-to-r from-purple-400 to-pink-400 opacity-0 group-hover:opacity-20 blur-xl transition-opacity" />
          </Button>
          
          <Button 
            size="lg" 
            variant="outline" 
            className="px-8 py-6 text-lg border-purple-500/30 bg-purple-900/20 hover:bg-purple-800/30 backdrop-blur-sm text-white"
          >
            <Play className="w-5 h-5 mr-2" />
            Tonton Demo
          </Button>
        </motion.div>

        {/* Stats */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1, delay: 1.2 }}
          className="mt-20 grid grid-cols-3 gap-8 max-w-2xl mx-auto"
        >
          {[
            { label: 'Video Diproses', value: 50, suffix: 'K+' },
            { label: 'Waktu Tersimpan', value: 10, suffix: 'Jt+ jam' },
            { label: 'Pengguna Aktif', value: 25, suffix: 'K+' }
          ].map((stat, i) => (
            <div key={i} className="text-center">
              <div className="text-3xl md:text-4xl font-bold text-white mb-2">
                <AnimatedCounter target={stat.value} />{stat.suffix}
              </div>
              <div className="text-sm text-purple-300">{stat.label}</div>
            </div>
          ))}
        </motion.div>
      </motion.div>

      {/* Nebula Orbs */}
      <div className="absolute top-1/4 -left-20 w-96 h-96 bg-purple-600/40 rounded-full blur-3xl" />
      <div className="absolute bottom-1/4 -right-20 w-96 h-96 bg-pink-600/40 rounded-full blur-3xl" />
      <div className="absolute top-1/2 left-1/2 w-64 h-64 bg-indigo-600/30 rounded-full blur-3xl" />
    </section>
  )
}

const FeaturesSection = () => {
  const features = [
    {
      icon: Brain,
      title: 'Analisis Bertenaga AI',
      description: 'AI canggih menganalisis konten Anda untuk mengidentifikasi momen paling menarik secara otomatis',
      color: 'from-purple-500 to-pink-500'
    },
    {
      icon: Scissors,
      title: 'Pemotongan Cerdas',
      description: 'Deteksi scene intelligent menciptakan klip dengan timing sempurna dan transisi halus',
      color: 'from-pink-500 to-rose-500'
    },
    {
      icon: Clock,
      title: 'Proses Instan',
      description: 'Proses konten berjam-jam dalam hitungan menit dengan pipeline AI kami yang dioptimalkan',
      color: 'from-orange-500 to-yellow-500'
    },
    {
      icon: Download,
      title: 'Multi Format',
      description: 'Export dalam format apapun yang dioptimalkan untuk YouTube Shorts, TikTok, Instagram, dan lainnya',
      color: 'from-blue-500 to-cyan-500'
    }
  ]

  return (
    <section className="relative py-32 overflow-hidden">
      <GalaxyBackground />
      
      <div className="container mx-auto px-6 relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="text-center mb-20"
        >
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/10 border border-purple-500/30 backdrop-blur-md mb-6">
            <Zap className="w-4 h-4 text-purple-300" />
            <span className="text-sm text-purple-200">Fitur</span>
          </div>
          <h2 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-b from-white to-purple-300 bg-clip-text text-transparent">
            Fitur AI yang Powerful
          </h2>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Semua yang Anda butuhkan untuk membuat klip viral dari konten long-form Anda
          </p>
        </motion.div>

        <div className="grid md:grid-cols-2 gap-8 max-w-6xl mx-auto">
          {features.map((feature, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
            >
              <Card className="group relative p-8 bg-gradient-to-br from-purple-900/30 to-indigo-900/20 border-purple-500/30 hover:border-purple-400/50 backdrop-blur-md transition-all duration-300 hover:scale-105 overflow-hidden">
                {/* Glow Effect */}
                <div className="absolute inset-0 bg-gradient-to-r opacity-0 group-hover:opacity-10 transition-opacity duration-300" 
                     style={{ backgroundImage: `linear-gradient(to right, var(--tw-gradient-stops))` }} />
                
                {/* Icon */}
                <div className="relative mb-6">
                  <div className={`inline-flex p-4 rounded-2xl bg-gradient-to-br ${feature.color} shadow-lg shadow-purple-500/50`}>
                    <feature.icon className="w-8 h-8 text-white" />
                  </div>
                </div>

                {/* Content */}
                <h3 className="text-2xl font-bold text-white mb-4">
                  {feature.title}
                </h3>
                <p className="text-gray-300 leading-relaxed">
                  {feature.description}
                </p>

                {/* Corner Accent */}
                <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-purple-400/10 to-transparent rounded-bl-full opacity-0 group-hover:opacity-100 transition-opacity" />
              </Card>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Nebula Effects */}
      <div className="absolute top-0 left-1/4 w-64 h-64 bg-purple-600/20 rounded-full blur-3xl" />
      <div className="absolute bottom-0 right-1/4 w-64 h-64 bg-pink-600/20 rounded-full blur-3xl" />
    </section>
  )
}

const HowItWorksSection = () => {
  const steps = [
    {
      number: '01',
      title: 'Paste URL Video',
      description: 'Cukup paste link video YouTube Anda',
      icon: Video
    },
    {
      number: '02',
      title: 'Analisis AI',
      description: 'AI kami menganalisis dan mengidentifikasi momen kunci',
      icon: Brain
    },
    {
      number: '03',
      title: 'Pemotongan Cerdas',
      description: 'Secara otomatis membuat klip yang dioptimalkan',
      icon: Scissors
    },
    {
      number: '04',
      title: 'Export & Bagikan',
      description: 'Download dan bagikan ke semua platform',
      icon: Download
    }
  ]

  return (
    <section className="relative py-32 bg-gradient-to-b from-black via-purple-950/20 to-black">
      <GalaxyBackground />
      
      <div className="container mx-auto px-6 relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="text-center mb-20"
        >
          <h2 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-b from-white to-purple-300 bg-clip-text text-transparent">
            Cara Kerjanya
          </h2>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Dari upload hingga klip viral hanya dalam empat langkah sederhana
          </p>
        </motion.div>

        <div className="max-w-5xl mx-auto">
          {steps.map((step, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, x: index % 2 === 0 ? -50 : 50 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: index * 0.2 }}
              className="relative mb-16 last:mb-0"
            >
              <div className={`flex flex-col md:flex-row items-center gap-8 ${index % 2 === 1 ? 'md:flex-row-reverse' : ''}`}>
                {/* Number & Icon */}
                <div className="relative flex-shrink-0">
                  <div className="relative w-32 h-32">
                    <div className="absolute inset-0 bg-gradient-to-br from-purple-500 to-pink-500 rounded-3xl transform rotate-6 opacity-30 blur-xl" />
                    <div className="relative w-full h-full bg-gradient-to-br from-purple-900/80 to-indigo-900/80 border border-purple-500/30 backdrop-blur-md rounded-3xl flex items-center justify-center">
                      <step.icon className="w-12 h-12 text-purple-300" />
                    </div>
                  </div>
                  <div className="absolute -top-4 -right-4 w-16 h-16 bg-gradient-to-br from-purple-500 to-pink-500 rounded-2xl flex items-center justify-center font-bold text-white shadow-lg shadow-purple-500/50">
                    {step.number}
                  </div>
                </div>

                {/* Content */}
                <div className={`flex-1 ${index % 2 === 1 ? 'md:text-right' : ''}`}>
                  <h3 className="text-3xl font-bold text-white mb-3">{step.title}</h3>
                  <p className="text-lg text-gray-300">{step.description}</p>
                </div>
              </div>

              {/* Connecting Line */}
              {index < steps.length - 1 && (
                <div className="hidden md:block absolute left-16 top-32 w-0.5 h-16 bg-gradient-to-b from-purple-500/50 to-transparent" />
              )}
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

const DemoSection = () => {
  const [isPlaying, setIsPlaying] = useState(false)

  return (
    <section className="relative py-32 overflow-hidden">
      <GalaxyBackground />
      
      <div className="container mx-auto px-6 relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="text-center mb-16"
        >
          <h2 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-b from-white to-purple-300 bg-clip-text text-transparent">
            Lihat Aksinya
          </h2>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Saksikan bagaimana AI kami mengubah video panjang menjadi klip yang engaging
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="max-w-5xl mx-auto"
        >
          <div className="relative group cursor-pointer" onClick={() => setIsPlaying(!isPlaying)}>
            {/* Video Preview */}
            <div className="relative aspect-video rounded-3xl overflow-hidden bg-gradient-to-br from-purple-900 to-indigo-900 border-2 border-purple-500/30 shadow-2xl shadow-purple-500/30">
              <img 
                src="https://images.unsplash.com/photo-1597733336794-12d05021d510?crop=entropy&cs=srgb&fm=jpg&ixid=M3w3NTY2NzR8MHwxfHNlYXJjaHwzfHxhYnN0cmFjdCUyMHRlY2h8ZW58MHx8fGJsdWV8MTc3MDg2Mjg2MXww&ixlib=rb-4.1.0&q=85"
                alt="Demo"
                className="w-full h-full object-cover opacity-70"
              />
              
              {/* Overlay */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-purple-900/40 to-transparent" />
              
              {/* Play Button */}
              <motion.div 
                className="absolute inset-0 flex items-center justify-center"
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
              >
                <div className="relative">
                  <div className="absolute inset-0 bg-purple-500 rounded-full blur-2xl opacity-50 animate-pulse" />
                  <div className="relative w-20 h-20 bg-white rounded-full flex items-center justify-center shadow-2xl">
                    <Play className="w-10 h-10 text-purple-600 ml-1" fill="currentColor" />
                  </div>
                </div>
              </motion.div>

              {/* Processing Indicator */}
              <div className="absolute bottom-6 left-6 right-6">
                <div className="backdrop-blur-md bg-purple-900/30 rounded-2xl p-4 border border-purple-500/30">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-white font-medium">AI Memproses Video...</span>
                    <span className="text-purple-300">85%</span>
                  </div>
                  <div className="h-2 bg-purple-950 rounded-full overflow-hidden">
                    <motion.div 
                      className="h-full bg-gradient-to-r from-purple-500 to-pink-500"
                      initial={{ width: '0%' }}
                      animate={{ width: '85%' }}
                      transition={{ duration: 2, ease: 'easeOut' }}
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Glow Effect */}
            <div className="absolute -inset-1 bg-gradient-to-r from-purple-500 to-pink-500 rounded-3xl opacity-0 group-hover:opacity-30 blur-xl transition-opacity" />
          </div>
        </motion.div>
      </div>

      {/* Nebula Accent */}
      <div className="absolute top-1/2 left-0 w-96 h-96 bg-purple-600/20 rounded-full blur-3xl" />
    </section>
  )
}

const PricingSection = () => {
  const plans = [
    {
      name: 'Pemula',
      price: 'Gratis',
      description: 'Sempurna untuk mencoba platform kami',
      features: [
        '5 video per bulan',
        'Video hingga 30 menit',
        'AI clipping dasar',
        'Export kualitas 720p'
      ],
      cta: 'Mulai Gratis',
      popular: false
    },
    {
      name: 'Pro',
      price: 'Rp299K',
      period: '/bulan',
      description: 'Untuk content creator dan profesional',
      features: [
        'Video unlimited',
        'Video hingga 3 jam',
        'AI clipping lanjutan',
        'Export kualitas 4K',
        'Proses prioritas',
        'Custom branding'
      ],
      cta: 'Mulai Sekarang',
      popular: true
    },
    {
      name: 'Enterprise',
      price: 'Custom',
      description: 'Untuk tim dan agensi',
      features: [
        'Semua fitur Pro',
        'Durasi video unlimited',
        'Akses API',
        'Support dedicated',
        'Training AI custom',
        'Solusi white-label'
      ],
      cta: 'Hubungi Sales',
      popular: false
    }
  ]

  return (
    <section className="relative py-32 bg-gradient-to-b from-black via-indigo-950/20 to-black">
      <GalaxyBackground />
      
      <div className="container mx-auto px-6 relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="text-center mb-20"
        >
          <h2 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-b from-white to-purple-300 bg-clip-text text-transparent">
            Harga Sederhana
          </h2>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Pilih paket yang sesuai kebutuhan Anda. Upgrade atau downgrade kapan saja.
          </p>
        </motion.div>

        <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          {plans.map((plan, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              className="relative"
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2 px-4 py-1 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full text-sm font-medium text-white shadow-lg shadow-purple-500/50">
                  Paling Populer
                </div>
              )}
              
              <Card className={`relative h-full p-8 bg-gradient-to-b from-purple-900/30 to-indigo-900/20 backdrop-blur-md transition-all duration-300 hover:scale-105 ${
                plan.popular ? 'border-purple-500/50 shadow-xl shadow-purple-500/30' : 'border-purple-500/20'
              }`}>
                {/* Plan Name */}
                <h3 className="text-2xl font-bold text-white mb-2">{plan.name}</h3>
                <p className="text-gray-300 text-sm mb-6">{plan.description}</p>

                {/* Price */}
                <div className="mb-8">
                  <span className="text-5xl font-bold text-white">{plan.price}</span>
                  {plan.period && <span className="text-gray-300 text-lg">{plan.period}</span>}
                </div>

                {/* Features */}
                <ul className="space-y-4 mb-8">
                  {plan.features.map((feature, i) => (
                    <li key={i} className="flex items-start gap-3">
                      <Check className="w-5 h-5 text-purple-400 flex-shrink-0 mt-0.5" />
                      <span className="text-gray-300">{feature}</span>
                    </li>
                  ))}
                </ul>

                {/* CTA Button */}
                <Button 
                  className={`w-full py-6 text-lg font-medium transition-all duration-300 ${
                    plan.popular 
                      ? 'bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 shadow-lg shadow-purple-500/50' 
                      : 'bg-purple-900/50 hover:bg-purple-800/50 border border-purple-500/30'
                  }`}
                >
                  {plan.cta}
                </Button>

                {/* Glow Effect for Popular */}
                {plan.popular && (
                  <div className="absolute -inset-1 bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg opacity-20 blur-xl -z-10" />
                )}
              </Card>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

const CTASection = () => {
  return (
    <section className="relative py-32 overflow-hidden">
      <GalaxyBackground />
      
      <div className="container mx-auto px-6 relative z-10">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="relative max-w-4xl mx-auto"
        >
          <div className="relative p-12 md:p-16 rounded-3xl bg-gradient-to-br from-purple-600 via-pink-500 to-indigo-600 overflow-hidden">
            {/* Animated Background Pattern */}
            <div className="absolute inset-0 opacity-20">
              <div className="absolute top-0 left-0 w-full h-full">
                {[...Array(20)].map((_, i) => (
                  <motion.div
                    key={i}
                    className="absolute w-2 h-2 bg-white rounded-full"
                    style={{
                      left: `${Math.random() * 100}%`,
                      top: `${Math.random() * 100}%`
                    }}
                    animate={{
                      opacity: [0.2, 0.8, 0.2],
                      scale: [1, 1.5, 1]
                    }}
                    transition={{
                      duration: Math.random() * 3 + 2,
                      repeat: Infinity,
                      ease: 'easeInOut'
                    }}
                  />
                ))}
              </div>
            </div>

            {/* Content */}
            <div className="relative z-10 text-center">
              <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">
                Siap Membuat Klip Viral?
              </h2>
              <p className="text-xl text-purple-100 mb-10 max-w-2xl mx-auto">
                Bergabung dengan ribuan kreator yang sudah menggunakan AI untuk mengembangkan audience mereka
              </p>
              
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button 
                  size="lg"
                  className="group px-8 py-6 text-lg bg-white text-purple-600 hover:bg-gray-100 border-0 shadow-xl"
                >
                  <span className="flex items-center gap-2">
                    Mulai Sekarang
                    <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                  </span>
                </Button>
                
                <Button 
                  size="lg"
                  variant="outline"
                  className="px-8 py-6 text-lg border-white/30 text-white hover:bg-white/10 backdrop-blur-sm"
                >
                  Jadwalkan Demo
                </Button>
              </div>
            </div>

            {/* Decorative Elements */}
            <div className="absolute -top-20 -left-20 w-40 h-40 bg-white/20 rounded-full blur-3xl" />
            <div className="absolute -bottom-20 -right-20 w-40 h-40 bg-pink-300/20 rounded-full blur-3xl" />
          </div>
        </motion.div>
      </div>

      {/* Background Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full max-w-4xl bg-purple-500/20 rounded-full blur-3xl" />
    </section>
  )
}

const Footer = () => {
  return (
    <footer className="relative border-t border-purple-900/30 py-12 bg-black/50 backdrop-blur-sm">
      <div className="container mx-auto px-6">
        <div className="grid md:grid-cols-4 gap-8 mb-8">
          <div>
            <h3 className="text-xl font-bold text-white mb-4">ClipifyAI</h3>
            <p className="text-gray-400 text-sm">
              Ubah video Anda menjadi klip viral dengan kekuatan AI
            </p>
          </div>
          
          <div>
            <h4 className="font-semibold text-white mb-4">Produk</h4>
            <ul className="space-y-2 text-sm text-gray-400">
              <li><a href="#" className="hover:text-purple-300 transition-colors">Fitur</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">Harga</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">API</a></li>
            </ul>
          </div>
          
          <div>
            <h4 className="font-semibold text-white mb-4">Perusahaan</h4>
            <ul className="space-y-2 text-sm text-gray-400">
              <li><a href="#" className="hover:text-purple-300 transition-colors">Tentang</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">Blog</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">Karir</a></li>
            </ul>
          </div>
          
          <div>
            <h4 className="font-semibold text-white mb-4">Legal</h4>
            <ul className="space-y-2 text-sm text-gray-400">
              <li><a href="#" className="hover:text-purple-300 transition-colors">Privasi</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">Ketentuan</a></li>
              <li><a href="#" className="hover:text-purple-300 transition-colors">Kontak</a></li>
            </ul>
          </div>
        </div>
        
        <div className="pt-8 border-t border-purple-900/30 text-center text-sm text-gray-400">
          <p>© 2026 ClipifyAI. Hak cipta dilindungi.</p>
        </div>
      </div>
    </footer>
  )
}

export default function App() {
  return (
    <div className="min-h-screen bg-black text-white overflow-x-hidden">
      <div className="relative">
        {/* Space Background */}
        <div className="fixed inset-0 bg-gradient-to-b from-indigo-950 via-purple-950 to-black -z-10" />
        
        <HeroSection />
        <FeaturesSection />
        <HowItWorksSection />
        <DemoSection />
        <PricingSection />
        <CTASection />
        <Footer />
      </div>
    </div>
  )
}
