package state_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-community/carousel/bosh"
	"github.com/cloudfoundry-community/carousel/credhub"
	. "github.com/cloudfoundry-community/carousel/state"
)

var _ = Describe("Credential", func() {
	var (
		olderThan     time.Time
		expiresBefore time.Time
		credential    *Credential
		criteria      RegenerationCriteria
	)

	BeforeEach(func() {
		olderThan = time.Now()
		expiresBefore = time.Now()
		criteria = RegenerationCriteria{
			OlderThan:     olderThan,
			ExpiresBefore: expiresBefore,
		}
	})

	Describe("NextAction", func() {
		BeforeEach(func() {
			vca := olderThan.Add(time.Hour)
			fooDeployment := Deployment{Name: "foo-deployment"}
			credential = &Credential{
				Latest:      true,
				Deployments: Deployments{&fooDeployment},
				Credential: &credhub.Credential{
					VersionCreatedAt: &vca,
					ID:               "foo-id",
					Name:             "/foo-name",
					Type:             credhub.Password,
				},
				Path: &Path{
					Deployments: Deployments{&fooDeployment},
				},
			}
			credential.Path.Versions = Credentials{credential}
		})

		Context("given a up-to-date credential", func() {
			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(None))
			})

		})

		Context("given a credential with its update mode set to no-overwrite", func() {
			BeforeEach(func() {
				credential.Path.VariableDefinition = &bosh.VariableDefinition{
					UpdateMode: bosh.NoOverwrite,
				}
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(NoOverwrite))
			})

			Context("but IgnoreUpdateMode was set to true", func() {
				BeforeEach(func() {
					criteria.IgnoreUpdateMode = true
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(None))
				})
			})

		})

		Context("given a latest credential which has not been deployed yet", func() {
			BeforeEach(func() {
				credential.Latest = true
				credential.Path.Deployments = append(
					credential.Path.Deployments, &Deployment{Name: "bar-deployment"})
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(BoshDeploy))
			})

			Context("which is too old", func() {
				BeforeEach(func() {
					vca := olderThan.Add(-10 * time.Minute)
					credential.VersionCreatedAt = &vca
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(Regenerate))
				})

				Context("of a non regeneratable type", func() {
					BeforeEach(func() {
						credential.Type = credhub.JSON
					})

					It("finds the next action", func() {
						Expect(credential.NextAction(criteria)).To(Equal(None))
					})
				})

			})
		})

		Context("given a cleanable credential which has not been deployed", func() {
			BeforeEach(func() {
				credential.Deployments = make(Deployments, 0)
				credential.Type = credhub.Certificate
				credential.Path.Deployments = make(Deployments, 0)
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(CleanUp))
			})

			Context("which is still being referenced by a deployed credential", func() {
				BeforeEach(func() {
					credential.ReferencedBy = Credentials{&Credential{
						Deployments: make(Deployments, 1),
					}}
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(None))
				})
			})

			Context("of a type which does not support version deletion", func() {
				BeforeEach(func() {
					credential.Type = credhub.Password
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(None))
				})
			})
		})

		Context("given a credential which is expiring", func() {
			BeforeEach(func() {
				ed := expiresBefore.Add(-time.Hour)
				credential.ExpiryDate = &ed
				edca := expiresBefore.Add(+time.Hour)
				credential.SignedBy = &Credential{
					Credential: &credhub.Credential{
						ExpiryDate: &edca,
					},
				}
				credential.SignedBy.Path = &Path{Versions: Credentials{credential.SignedBy}}
			})

			It("finds the next action", func() {
				credential.PathVersion()
				Expect(credential.NextAction(criteria)).To(Equal(Regenerate))
			})

			Context("which is self-signed", func() {
				BeforeEach(func() {
					credential.SignedBy = nil
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(Regenerate))
				})
			})

			Context("which is signed by an expiring ca", func() {
				BeforeEach(func() {
					ed := expiresBefore.Add(-time.Hour)
					credential.SignedBy = &Credential{
						Credential: &credhub.Credential{
							ExpiryDate: &ed,
						},
					}
					credential.SignedBy.Path = &Path{
						Versions: Credentials{credential.SignedBy},
					}
				})

				It("finds the next action", func() {
					Expect(credential.NextAction(criteria)).To(Equal(None))
				})

				Context("with a active not expiring sibiling ca", func() {
					BeforeEach(func() {
						ed := expiresBefore.Add(+time.Hour)
						credential.SignedBy.Path.Versions = append(
							credential.SignedBy.Path.Versions,
							&Credential{Latest: true,
								ReferencedBy: Credentials{credential},
								Credential: &credhub.Credential{
									Transitional: false,
									ExpiryDate:   &ed,
								}},
						)
					})

					It("finds the next action", func() {
						Expect(credential.NextAction(criteria)).To(Equal(Regenerate))
					})
				})
			})
		})

		Context("given a valid leaf with a singing ca that has an active non transitional latest sibling", func() {
			BeforeEach(func() {
				signing := true
				credential.Latest = true
				credential.SignedBy = &Credential{
					Latest: false,
					Credential: &credhub.Credential{
						Transitional: true,
					},
					Path: &Path{},
				}
				activeCa := &Credential{
					Latest:       true,
					Signing:      &signing,
					ReferencedBy: Credentials{credential},
					Credential: &credhub.Credential{
						Transitional: false,
					},
				}
				credential.SignedBy.Path.Versions = Credentials{activeCa, credential.SignedBy}
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(Regenerate))
			})
		})

		Context("given a transitional latest credential which is not referenced", func() {
			BeforeEach(func() {
				credential.Latest = true
				credential.Transitional = true
				credential.SignedBy = nil
				credential.Type = credhub.Certificate
				credential.ReferencedBy = make(Credentials, 0)
				credential.Deployments = make(Deployments, 0)
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(CleanUp))
			})
		})

		Context("given a signing credential with a latest active transitional sibling", func() {
			BeforeEach(func() {
				signing := true
				vca := olderThan.Add(time.Hour)
				credential.Signing = &signing
				credential.Latest = false
				credential.Path.Versions = append(credential.Path.Versions,
					&Credential{
						Deployments: credential.Path.Deployments,
						Latest:      true,
						Credential: &credhub.Credential{
							VersionCreatedAt: &vca,
							Transitional:     true,
						},
						Path: credential.Path,
					})
			})

			It("finds the next action", func() {
				Expect(credential.NextAction(criteria)).To(Equal(MarkTransitional))
			})
		})

		Context("given a transitional credential with an latest active sibling", func() {
			Context("which signs all credentials signed by self", func() {
				Context("all of which have also been deployed", func() {
					BeforeEach(func() {
						signing := true
						vca := olderThan.Add(time.Hour)
						signingCa := Credential{
							Deployments: credential.Path.Deployments,
							Latest:      true,
							Signing:     &signing,
							Credential: &credhub.Credential{
								VersionCreatedAt: &vca,
							},
							Path: credential.Path,
						}

						oldLeaf := Credential{
							Deployments: make(Deployments, 0),
							SignedBy:    credential,
						}

						newLeaf := Credential{
							Deployments: credential.Path.Deployments,
							SignedBy:    &signingCa,
						}

						leafPath := Path{
							Versions: Credentials{&newLeaf, &oldLeaf},
						}

						oldLeaf.Path = &leafPath
						newLeaf.Path = &leafPath

						credential.Transitional = true
						credential.Path.Versions = append(credential.Path.Versions, &signingCa)

						credential.Signs = Credentials{&oldLeaf}
						signingCa.Signs = Credentials{&newLeaf}

						credential.Deployments = make(Deployments, 0)
						credential.Latest = false
					})

					It("finds the next action", func() {
						Expect(credential.NextAction(criteria)).To(Equal(UnMarkTransitional))
					})
				})
			})
		})
	})
})
